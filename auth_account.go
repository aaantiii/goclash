package goclash

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/bytedance/sonic"
)

func (h *AccountsKeyProvider) RevalidateKeys() error {
	if err := h.updateIPAddr(); err != nil {
		return err
	}
	return h.updateAccounts()
}

func (h *AccountsKeyProvider) updateIPAddr() error {
	res, err := h.rc.R().Get(IPifyEndpoint)
	if err != nil {
		return err
	}

	body := string(res.Body())
	if res.StatusCode() != http.StatusOK {
		return errors.New(body)
	}
	if body == "" {
		return errors.New("couldn't get IP address")
	}
	if body == h.ipAddr {
		return nil
	}

	h.mu.Lock()
	h.ipAddr = body
	h.mu.Unlock()
	return nil
}

func (h *AccountsKeyProvider) updateAccounts() error {
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

func (h *AccountsKeyProvider) login(account *APIAccount) error {
	res, err := h.rc.NewDefaultRequest().SetBody(account.Credentials).Post(DevLoginEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	return sonic.Unmarshal(res.Body(), &account)
}

// getAccountKeys retrieves the API keys for the given account and sets APIAccount.Keys.
func (h *AccountsKeyProvider) getAccountKeys(account *APIAccount) error {
	res, err := h.rc.NewDefaultRequest().Post(DevKeyListEndpoint.URL())
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
	for i := range body.Keys {
		account.Keys[i] = body.Keys[i]
	}
	h.mu.Unlock()
	return nil
}

func (h *AccountsKeyProvider) updateAccountKeys(account *APIAccount) error {
	if err := h.getAccountKeys(account); err != nil {
		return err
	}

	errChan := make(chan error, keysPerAccount)
	var freeKeyIndexes []int
	var wg sync.WaitGroup
	for i := 0; i < keysPerAccount; i++ {
		if account.Keys[i] == nil {
			freeKeyIndexes = append(freeKeyIndexes, i)
			continue
		}
		if !slices.Contains(account.Keys[i].CidrRanges, h.ipAddr) {
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
			}(account.Keys[i], i)
		}
	}
	wg.Wait()

	for i := range freeKeyIndexes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := h.createAccountKey(account, i); err != nil {
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

func (h *AccountsKeyProvider) createAccountKey(account *APIAccount, index int) error {
	desc := fmt.Sprintf("Created at %s by goclash", time.Now().UTC().Round(time.Minute).String())
	key := &APIKey{
		Name:        "goclash",
		Description: desc,
		CidrRanges:  []string{h.ipAddr},
		Scopes:      []string{"clash"},
	}
	res, err := h.rc.NewDefaultRequest().SetBody(key).Post(DevKeyCreateEndpoint.URL())
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

func (h *AccountsKeyProvider) revokeAccountKey(key *APIKey) error {
	payload := map[string]string{"id": key.ID}
	res, err := h.rc.NewDefaultRequest().SetBody(payload).Post(DevKeyRevokeEndpoint.URL())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}
	return nil
}

func (h *AccountsKeyProvider) GetKey() string {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := h.accounts[h.keyIndex.AccountIndex].Keys[h.keyIndex.KeyIndex]
	if h.keyIndex.KeyIndex == len(h.accounts[h.keyIndex.AccountIndex].Keys)-1 {
		h.keyIndex.AccountIndex = (h.keyIndex.AccountIndex + 1) % len(h.accounts)
	}
	h.keyIndex.KeyIndex = (h.keyIndex.KeyIndex + 1) % len(h.accounts[h.keyIndex.AccountIndex].Keys)
	return key.Key
}
