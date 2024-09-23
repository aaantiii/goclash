package goclash

import (
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type KeyProvider interface {
	GetKey() string
	// Does all the necessary work to update the key provider, making the keys valid ones again
	RevalidateKeys() error
}

type Client struct {
	keys  KeyProvider
	rc    *RequestsClient
	cache *Cache
}

func newClient(creds Credentials) (*Client, error) {
	accounts := make([]*APIAccount, 0, len(creds))
	for email, password := range creds {
		accounts = append(accounts, &APIAccount{
			Credentials: &APIAccountCredentials{
				Email:    email,
				Password: password,
			},
		})
	}

	rc := &RequestsClient{resty.New()}
	keys := &AccountsKeyProvider{
		accounts: accounts,
		rc:       rc,
	}
	client := &Client{
		keys:  keys,
		rc:    rc,
		cache: newCache(),
	}

	if err := keys.RevalidateKeys(); err != nil {
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
			if err = h.keys.RevalidateKeys(); err != nil {
				return nil, err
			}
			return h.do(method, url, req, false)
		}
	}

	return nil, clientErr
}

func (h *Client) newDefaultRequest() *resty.Request {
	return h.rc.NewDefaultRequest()
}

func (h *Client) withAuth(req *resty.Request) *resty.Request {
	return req.SetAuthToken(h.keys.GetKey())
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
