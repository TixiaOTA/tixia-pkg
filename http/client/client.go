package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries        int           // Maximum number of retry attempts (default: 3)
	InitialBackoff    time.Duration // Initial backoff duration (default: 100ms)
	MaxBackoff        time.Duration // Maximum backoff duration (default: 5s)
	BackoffMultiplier float64       // Multiplier for exponential backoff (default: 2.0)
}

type HTTPConfig struct {
	BaseURL  string
	Timeout  *time.Duration // Timeout for individual request attempts
	Deadline *time.Duration // Deadline for the entire operation (including retries)
	Retry    *RetryConfig   // Optional retry configuration
}

type HTTPClient struct {
	config HTTPConfig
}

func NewHTTPClient(config HTTPConfig) *HTTPClient {
	// Set default retry config if not provided
	if config.Retry == nil {
		config.Retry = &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        5 * time.Second,
			BackoffMultiplier: 2.0,
		}
	}

	return &HTTPClient{
		config: config,
	}
}

// SetBaseURL sets the base URL for the HTTP client
func (c *HTTPClient) SetBaseURL(baseURL string) {
	c.config.BaseURL = baseURL
}

// SetTimeout sets the timeout for individual request attempts
func (c *HTTPClient) SetTimeout(timeout time.Duration) {
	c.config.Timeout = &timeout
}

// SetDeadline sets the deadline for the entire operation (including retries)
func (c *HTTPClient) SetDeadline(deadline time.Duration) {
	c.config.Deadline = &deadline
}

// SetRetry sets the retry configuration for the HTTP client
func (c *HTTPClient) SetRetry(retry RetryConfig) {
	c.config.Retry = &retry
}

// SetMaxRetries sets the maximum number of retries for the HTTP client
func (c *HTTPClient) SetMaxRetries(maxRetries int) {
	c.config.Retry.MaxRetries = maxRetries
}

// SetInitialBackoff sets the initial backoff for the HTTP client
func (c *HTTPClient) SetInitialBackoff(initialBackoff time.Duration) {
	c.config.Retry.InitialBackoff = initialBackoff
}

// SetMaxBackoff sets the maximum backoff for the HTTP client
func (c *HTTPClient) SetMaxBackoff(maxBackoff time.Duration) {
	c.config.Retry.MaxBackoff = maxBackoff
}

// SetBackoffMultiplier sets the backoff multiplier for the HTTP client
func (c *HTTPClient) SetBackoffMultiplier(backoffMultiplier float64) {
	c.config.Retry.BackoffMultiplier = backoffMultiplier
}

// Get makes a GET request to the given path. Base URL is prepended to the path.
func (c *HTTPClient) Get(ctx context.Context, request GetRequest) (*http.Response, error) {
	var timeout time.Duration

	// Get timeout from config
	if c.config.Timeout == nil {
		timeout = 30 * time.Second // Default timeout
	} else {
		timeout = *c.config.Timeout
	}

	// Create deadline context if configured
	deadlineCtx := ctx
	if c.config.Deadline != nil {
		var cancel context.CancelFunc
		deadlineCtx, cancel = context.WithTimeout(ctx, *c.config.Deadline)
		defer cancel()
	}

	// Create request with a temporary context (will be replaced per attempt)
	req, err := http.NewRequestWithContext(deadlineCtx, http.MethodGet, c.config.BaseURL+request.Path, nil)
	if err != nil {
		return nil, err
	}

	// Add additional headers
	if request.Headers != nil {
		for key, value := range *request.Headers {
			req.Header.Set(key, value)
		}
	}

	// Send request with retry
	resp, err := c.executeWithRetry(deadlineCtx, req, timeout)
	if err != nil {
		return nil, err
	}

	// Return response
	return resp, nil
}

// PostJSON makes a POST request with JSON body to the given path. Base URL is prepended to the path.
func (c *HTTPClient) PostJSON(ctx context.Context, request PostJSONRequest) (*http.Response, error) {
	var timeout time.Duration

	// Get timeout from config
	if c.config.Timeout == nil {
		timeout = 30 * time.Second // Default timeout
	} else {
		timeout = *c.config.Timeout
	}

	// Create deadline context if configured
	deadlineCtx := ctx
	if c.config.Deadline != nil {
		var cancel context.CancelFunc
		deadlineCtx, cancel = context.WithTimeout(ctx, *c.config.Deadline)
		defer cancel()
	}

	// Create request
	jsonBody, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequestWithContext(deadlineCtx, http.MethodPost, c.config.BaseURL+request.Path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set GetBody to allow body to be re-read for retries
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(jsonBody)), nil
	}

	// Add content type header
	req.Header.Set("Content-Type", "application/json")

	// Add additional headers
	if request.Headers != nil {
		for key, value := range *request.Headers {
			req.Header.Set(key, value)
		}
	}

	// Send request with retry
	resp, err := c.executeWithRetry(deadlineCtx, req, timeout)
	if err != nil {
		return nil, err
	}

	// Return response
	return resp, nil
}

// PostForm makes a POST request with form body to the given path. Base URL is prepended to the path.
func (c *HTTPClient) PostForm(ctx context.Context, request PostFormRequest) (*http.Response, error) {
	var timeout time.Duration

	// Get timeout from config
	if c.config.Timeout == nil {
		timeout = 30 * time.Second // Default timeout
	} else {
		timeout = *c.config.Timeout
	}

	// Create deadline context if configured
	deadlineCtx := ctx
	if c.config.Deadline != nil {
		var cancel context.CancelFunc
		deadlineCtx, cancel = context.WithTimeout(ctx, *c.config.Deadline)
		defer cancel()
	}

	// Create request
	formBodyBytes := []byte(request.Body.Encode())
	bodyReader := bytes.NewReader(formBodyBytes)
	req, err := http.NewRequestWithContext(deadlineCtx, http.MethodPost, c.config.BaseURL+request.Path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set GetBody to allow body to be re-read for retries
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(formBodyBytes)), nil
	}

	// Add content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add additional headers
	if request.Headers != nil {
		for key, value := range *request.Headers {
			req.Header.Set(key, value)
		}
	}

	// Send request with retry
	resp, err := c.executeWithRetry(deadlineCtx, req, timeout)
	if err != nil {
		return nil, err
	}

	// Return response
	return resp, nil
}
