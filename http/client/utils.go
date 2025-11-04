package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// executeWithRetry executes an HTTP request with retry logic
func (c *HTTPClient) executeWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response
	backoff := c.config.Retry.InitialBackoff

	for attempt := 0; attempt <= c.config.Retry.MaxRetries; attempt++ {
		// Execute the request
		resp, err := http.DefaultClient.Do(req)

		// Check if we should retry
		if !shouldRetry(err, resp) {
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		// Store the error/response for potential return
		lastErr = err

		// If we got a response that needs retry, close it before retrying
		// (unless this is the last attempt, then we'll return it)
		if resp != nil && attempt < c.config.Retry.MaxRetries {
			// Close the response body to free resources before retry
			resp.Body.Close()
		}

		if resp != nil {
			lastResp = resp
		}

		// If this was the last attempt, break
		if attempt == c.config.Retry.MaxRetries {
			break
		}

		// Wait before retrying (exponential backoff)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}

		// Increase backoff for next retry
		backoff = time.Duration(float64(backoff) * c.config.Retry.BackoffMultiplier)
		if backoff > c.config.Retry.MaxBackoff {
			backoff = c.config.Retry.MaxBackoff
		}

		// Recreate request body if it was consumed
		if req.Body != nil {
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				req.Body = body
			}
		}
	}

	// Return the last error or response
	if lastErr != nil {
		return lastResp, lastErr
	}
	return lastResp, nil
}

// shouldRetry determines if an error or response status should trigger a retry
func shouldRetry(err error, resp *http.Response) bool {
	// Retry on network errors
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			// Retry on timeout or connection errors
			if urlErr.Timeout() || urlErr.Temporary() {
				return true
			}
		}
		// Retry on context deadline exceeded
		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}
		return false
	}

	// Retry on HTTP 5xx status codes
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}

	return false
}
