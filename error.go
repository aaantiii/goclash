package clash

import (
	"fmt"
	"net/http"
)

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

func newClientErr(err error) *ClientError {
	return &ClientError{
		Status: http.StatusInternalServerError,
		APIError: &APIError{
			Reason:  ReasonBadRequest,
			Message: fmt.Sprintf("an unknown error occured: %s", err.Error()),
			Type:    "clash.go",
		},
	}
}

func (e *ClientError) Error() string {
	return e.Message
}
