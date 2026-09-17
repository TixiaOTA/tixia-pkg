package panorama

import (
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

// RequestOption configures one Panorama request.
type RequestOption func(*panoramaRequestOptions)

// SellJourneysOption is kept as an alias for backwards compatibility.
type SellJourneysOption = RequestOption

type panoramaRequestOptions struct {
	timeout time.Duration
	retry   *client.RetryConfig
}

// WithTimeout overrides the timeout for one request attempt.
func WithTimeout(timeout time.Duration) RequestOption {
	return func(options *panoramaRequestOptions) {
		options.timeout = timeout
	}
}

// WithRetry overrides the retry configuration for one request.
func WithRetry(retry client.RetryConfig) RequestOption {
	return func(options *panoramaRequestOptions) {
		options.retry = &retry
	}
}

func buildRequestOptions(defaultTimeout time.Duration, options ...RequestOption) client.RequestOptions {
	requestOptions := panoramaRequestOptions{timeout: defaultTimeout}
	for _, option := range options {
		if option != nil {
			option(&requestOptions)
		}
	}

	return client.RequestOptions{
		Timeout: &requestOptions.timeout,
		Retry:   requestOptions.retry,
	}
}
