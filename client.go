package goclash

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	accounts []*APIAccount
	rc       *resty.Client
	ipAddr   string
	keyIndex APIKeyIndex
	cache    *Cache
	mu       sync.Mutex
}

var defaultHeaders = map[string]string{
	"Accept":       "application/json",
	"Content-Type": "application/json",
	"User-Agent":   "goclash",
}

func newClient(creds Credentials) (*Client, error) {
	accounts := make([]*APIAccount, 0, len(creds))
	for email, password := range creds {
		accounts = append(accounts, &APIAccount{
			Keys: make([]*APIKey, keysPerAccount),
			Credentials: &APIAccountCredentials{
				Email:    email,
				Password: password,
			},
		})
	}

	client := &Client{
		accounts: accounts,
		rc:       resty.New(),
		cache:    newCache(),
	}

	if err := client.updateIPAddr(); err != nil {
		return nil, err
	}
	if err := client.updateAccounts(); err != nil {
		return nil, err
	}

	return client, nil
}

func (h *Client) do(method, url string, req *resty.Request, retry bool) ([]byte, error) {
	if h.cache.enabled {
		if data, ok := h.cache.Get(url); ok {
			return data, nil
		}
	}

	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() < 300 {
		h.cache.CacheResponse(url, res)
		return res.Body(), nil
	}

	clientErr := &ClientError{Status: res.StatusCode(), APIError: &APIError{}}
	if err = sonic.Unmarshal(res.Body(), &clientErr.APIError); err != nil {
		return nil, err
	}
	if res.StatusCode() == http.StatusForbidden {
		if !retry {
			return nil, clientErr
		}

		if clientErr.APIError.Reason == ReasonInvalidIP {
			if err = h.updateIPAddr(); err != nil {
				return nil, err
			}
			if err = h.updateAccounts(); err != nil {
				return nil, err
			}
			return h.do(method, url, req, false)
		}
	}

	return nil, clientErr
}

func (h *Client) updateIPAddr() error {
	res, err := h.rc.R().Get(IPifyEndpoint)
	if err != nil {
		return err
	}

	body := string(res.Body())
	if res.StatusCode() != http.StatusOK {
		return errors.New(body)
	}
	if body == h.ipAddr {
		return nil
	}

	h.mu.Lock()
	h.ipAddr = body
	h.mu.Unlock()
	return nil
}

func (h *Client) login(account *APIAccount) error {
	res, err := h.newDefaultRequest().SetBody(account.Credentials).Post(DevLoginEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	return sonic.Unmarshal(res.Body(), &account)
}

func (h *Client) updateAccounts() error {
	for _, account := range h.accounts {
		if err := h.login(account); err != nil {
			return err
		}
		if err := h.updateAccountKeys(account); err != nil {
			return err
		}
	}

	return nil
}

// getAccountKeys retrieves the API keys for the given account and sets APIAccount.Keys.
func (h *Client) getAccountKeys(account *APIAccount) error {
	res, err := h.newDefaultRequest().Post(DevKeyListEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	var body *KeyListResponse
	if err = sonic.Unmarshal(res.Body(), &body); err != nil {
		return err
	}

	h.mu.Lock()
	account.Keys = body.Keys
	h.mu.Unlock()
	return nil
}

func (h *Client) updateAccountKeys(account *APIAccount) error {
	if err := h.getAccountKeys(account); err != nil {
		return err
	}

	errChan := make(chan error, keysPerAccount)
	var freeKeyIndexes []int
	var wg sync.WaitGroup
	for i, key := range account.Keys {
		if key == nil {
			freeKeyIndexes = append(freeKeyIndexes, i)
			continue
		}
		if !slices.Contains(key.CidrRanges, h.ipAddr) {
			wg.Add(1)
			go func(key *APIKey, i int) {
				defer wg.Done()
				if err := h.revokeAccountKey(key); err != nil {
					errChan <- err
					return
				}
				if err := h.createAccountKey(account, i); err != nil {
					errChan <- err
					return
				}
			}(key, i)
		}
	}
	wg.Wait()

	for i := range freeKeyIndexes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := h.createAccountKey(account, keysPerAccount-i); err != nil {
				errChan <- err
			}
		}(i)
	}
	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (h *Client) createAccountKey(account *APIAccount, index int) error {
	desc := fmt.Sprintf("Created at %s by goclash", time.Now().UTC().Round(time.Minute).String())
	res, err := h.newDefaultRequest().SetBody(&APIKey{
		Name:        "goclash",
		Description: desc,
		CidrRanges:  []string{h.ipAddr},
		Scopes:      []string{"clash"},
	}).Post(DevKeyCreateEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	var keyRes *CreateKeyResponse
	if err = sonic.Unmarshal(res.Body(), &keyRes); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	account.Keys[index] = keyRes.Key
	return nil
}

func (h *Client) revokeAccountKey(key *APIKey) error {
	payload := map[string]any{"id": key}
	res, err := h.newDefaultRequest().SetBody(payload).Post(DevKeyRevokeEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}
	return nil
}

func (h *Client) getKey() string {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := h.accounts[h.keyIndex.AccountIndex].Keys[h.keyIndex.KeyIndex]
	if h.keyIndex.KeyIndex == len(h.accounts[h.keyIndex.AccountIndex].Keys)-1 {
		h.keyIndex.AccountIndex = (h.keyIndex.AccountIndex + 1) % len(h.accounts)
	}
	h.keyIndex.KeyIndex = (h.keyIndex.KeyIndex + 1) % len(h.accounts[h.keyIndex.AccountIndex].Keys)
	return key.Key
}

func (h *Client) newDefaultRequest() *resty.Request {
	return h.rc.R().SetHeaders(defaultHeaders)
}

func (h *Client) withAuth(req *resty.Request) *resty.Request {
	return req.SetAuthToken(h.getKey())
}

func (h *Client) withPaging(r *resty.Request, params *PagingParams) *resty.Request {
	if params == nil {
		return r
	}

	if params.After != "" {
		r.SetQueryParam("after", params.After)
	} else if params.Before != "" {
		r.SetQueryParam("before", params.Before)
	}
	if params.Limit > 0 {
		r.SetQueryParam("limit", strconv.Itoa(params.Limit))
	}
	return r
}
