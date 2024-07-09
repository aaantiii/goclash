package goclash

const keysPerAccount = 10

type APIAccount struct {
	Credentials *APIAccountCredentials
	Keys        [keysPerAccount]*APIKey
}

type APIAccountCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Credentials is a map of email to password.
type Credentials map[string]string

type APIKey struct {
	ID          string   `json:"id"`
	Origins     any      `json:"origins"`
	ValidUntil  any      `json:"validUntil"`
	DeveloperID string   `json:"developerId"`
	Tier        string   `json:"tier"`
	Name        string   `json:"name"`
	Description any      `json:"description"`
	Key         string   `json:"key"`
	Scopes      []string `json:"scopes"`
	CidrRanges  []string `json:"cidrRanges"`
}

// APIKeyIndex is used to determine which account and key to use for a given request.
type APIKeyIndex struct {
	AccountIndex int
	KeyIndex     int
}

type Developer struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Game          string   `json:"game"`
	Tier          string   `json:"tier"`
	AllowedScopes []string `json:"allowedScopes"`
	MaxCidrs      int      `json:"maxCidrs"`
	PrevLoginTS   string   `json:"prevLoginTs"`
	PrevLoginIP   string   `json:"prevLoginIp"`
	PrevLoginUA   string   `json:"prevLoginUa"`
}

type CreateKeyResponse struct {
	Key                     *APIKey `json:"key,omitempty"`
	Status                  Status  `json:"status"`
	SessionExpiresInSeconds int     `json:"sessionExpiresInSeconds"`
}

type KeyListResponse struct {
	Keys                    []*APIKey `json:"keys,omitempty"`
	Status                  Status    `json:"status"`
	SessionExpiresInSeconds int       `json:"sessionExpiresInSeconds"`
}

type Status struct {
	Detail  any    `json:"detail"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
