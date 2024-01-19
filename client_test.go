package clash

import (
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var client = newTestClient()

func newTestClient() *Client {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	creds := make(Credentials)
	emailStr := os.Getenv("EMAILS")
	passwordStr := os.Getenv("PASSWORDS")
	emails := strings.Split(emailStr, ",")
	passwords := strings.Split(passwordStr, ",")

	for i, email := range emails {
		creds[email] = passwords[i]
	}

	testClient, err := newClient(creds)
	if err != nil {
		panic(err)
	}

	return testClient
}

func Test_client(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{{
		name:    "Test Account 1",
		wantErr: false,
	}, {
		name:    "Test Account 2",
		wantErr: false,
	}}

	for i, tt := range tests {
		err := client.updateAccountKeys(client.accounts[i])
		if (err != nil) != tt.wantErr {
			t.Errorf("client.createAccountKey() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
