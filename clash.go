package clash

func NewClient(creds Credentials) (*Client, error) {
	return newClient(creds)
}
