package clash

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	accounts []*APIAccount
	rc       *resty.Client
	ipAddr   string
	keyIndex APITokenIndexer
	mu       sync.Mutex
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
		rc:       resty.New(),
		accounts: accounts,
	}

	if err := client.updateIPAddr(); err != nil {
		return nil, err
	}

	if err := client.updateAccounts(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) do(method string, route Endpoint, body any, retry bool) ([]byte, error) {
	req := c.withAuth(c.newDefaultRequest())
	if method == http.MethodPost && body != nil {
		req.SetBody(body)
	}

	res, err := req.Execute(method, route.URL())
	if err != nil {
		return nil, err
	}

	clientErr := &ClientError{Status: res.StatusCode()}
	if clientErr.Status >= 300 {
		if err = sonic.Unmarshal(res.Body(), &clientErr.APIError); err != nil {
			return nil, err
		}
	}

	if res.StatusCode() == http.StatusForbidden {
		if !retry {
			return nil, clientErr
		}

		if clientErr.APIError.Reason == ReasonInvalidIP {
			if err = c.updateIPAddr(); err != nil {
				return nil, err
			}
			if err = c.updateAccounts(); err != nil {
				return nil, err
			}
			return c.do(method, route, body, false)
		}
	}

	return res.Body(), nil
}

func (c *Client) withAuth(req *resty.Request) *resty.Request {
	key := c.getKey()
	return req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", key))
}

func (c *Client) updateIPAddr() error {
	res, err := c.rc.R().Get(IPifyEndpoint)
	if err != nil {
		return err
	}

	body := string(res.Body())
	if res.StatusCode() != http.StatusOK {
		return errors.New(body)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.ipAddr = body
	return nil
}

func (c *Client) login(account *APIAccount) error {
	res, err := c.newDefaultRequest().SetBody(account.Credentials).Post(DevLoginEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	return sonic.Unmarshal(res.Body(), &account)
}

func (c *Client) updateAccounts() error {
	for _, account := range c.accounts {
		if err := c.login(account); err != nil {
			return err
		}
		if err := c.updateAccountKeys(account); err != nil {
			return err
		}
	}

	return nil
}

// getAccountKeys retrieves the API keys for the given account and sets APIAccount.Keys.
func (c *Client) getAccountKeys(account *APIAccount) error {
	res, err := c.newDefaultRequest().Post(DevKeyListEndpoint.URL())
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

	c.mu.Lock()
	defer c.mu.Unlock()
	account.Keys = body.Keys
	return nil
}

func (c *Client) createAccountKey(account *APIAccount, index int) error {
	desc := fmt.Sprintf("Created at %s by clash.go", time.Now().UTC().Round(time.Minute).String())
	res, err := c.newDefaultRequest().SetBody(&APIKey{
		Name:        "clash.go",
		Description: desc,
		CidrRanges:  []string{c.ipAddr},
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

	c.mu.Lock()
	defer c.mu.Unlock()
	account.Keys[index] = keyRes.Key
	return nil
}

func (c *Client) updateAccountKeys(account *APIAccount) error {
	if err := c.getAccountKeys(account); err != nil {
		return err
	}

	errChan := make(chan error, keysPerAccount)
	var freeKeysIndex []int
	var wg sync.WaitGroup
	for i, key := range account.Keys {
		if key == nil {
			freeKeysIndex = append(freeKeysIndex, i)
			continue
		}
		if !slices.Contains(key.CidrRanges, c.ipAddr) {
			wg.Add(1)
			go func(key *APIKey, i int) {
				defer wg.Done()
				if err := c.revokeAccountKey(key); err != nil {
					errChan <- err
					return
				}
				if err := c.createAccountKey(account, i); err != nil {
					errChan <- err
					return
				}
			}(key, i)
		}
	}
	wg.Wait()

	for i := range freeKeysIndex {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := c.createAccountKey(account, keysPerAccount-i)
			if err != nil {
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

func (c *Client) revokeAccountKey(key *APIKey) error {
	payload := struct {
		ID string `json:"id"`
	}{ID: key.ID}
	res, err := c.newDefaultRequest().SetBody(payload).Post(DevKeyRevokeEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	return nil
}

func (c *Client) getKey() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.accounts[c.keyIndex.AccountIndex].Keys[c.keyIndex.TokenIndex]
	c.keyIndex.TokenIndex = (c.keyIndex.TokenIndex + 1) % len(c.accounts[c.keyIndex.AccountIndex].Keys)
	if c.keyIndex.TokenIndex == len(c.accounts[c.keyIndex.AccountIndex].Keys)-1 {
		c.keyIndex.AccountIndex = (c.keyIndex.AccountIndex + 1) % len(c.accounts)
	}

	return key.Key
}

func (c *Client) newDefaultRequest() *resty.Request {
	return c.rc.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "clash.go")
}
