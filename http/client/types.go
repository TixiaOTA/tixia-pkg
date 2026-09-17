package client

import (
	"net/url"
	"time"
)

// RequestOptions overrides HTTP settings for one request.
type RequestOptions struct {
	Timeout *time.Duration
	Retry   *RetryConfig
}

// GetRequest is the request for the Get method
type GetRequest struct {
	Path    string
	Headers *map[string]string // Optional additional headers
	Options *RequestOptions
}

// PostRequest is the request for the Post method
type PostJSONRequest struct {
	Path    string
	Body    any                // JSON body
	Headers *map[string]string // Optional additional headers
	Options *RequestOptions
}

// PostFormRequest is the request for the PostForm method
type PostFormRequest struct {
	Path    string
	Body    url.Values         // Form body
	Headers *map[string]string // Optional additional headers
	Options *RequestOptions
}
