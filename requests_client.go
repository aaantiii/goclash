package goclash

import "github.com/go-resty/resty/v2"

type RequestsClient struct {
	*resty.Client
}

var defaultHeaders = map[string]string{
	"Accept":       "application/json",
	"Content-Type": "application/json",
	"User-Agent":   "goclash",
}

func (rc *RequestsClient) NewDefaultRequest() *resty.Request {
	return rc.R().SetHeaders(defaultHeaders)
}
