package goclash

// APIError is the error directly returned by the Clash of Clans API. Every error returned by Client is ClientError, which embeds *APIError.
type APIError struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

const (
	ReasonBadRequest           = "badRequest"
	ReasonInvalidAuthorization = "accessDenied"
	ReasonInvalidIP            = "accessDenied.invalidIp"
	ReasonNotFound             = "notFound"
)

// ClientError is the error type returned by the client.
type ClientError struct {
	*APIError
	Status int `json:"status"`
}

func (e *ClientError) Error() string {
	return e.Message
}
