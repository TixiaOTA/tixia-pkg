package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// executeWithRetry executes an HTTP request with retry logic
// deadlineCtx: Context for overall operation deadline (stops retries when deadline is reached)
// req: HTTP request to execute
// timeout: Timeout duration for individual request attempts
func (c *HTTPClient) executeWithRetry(deadlineCtx context.Context, req *http.Request, timeout time.Duration) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response
	backoff := c.config.Retry.InitialBackoff

	for attempt := 0; attempt <= c.config.Retry.MaxRetries; attempt++ {
		// Check if deadline has been reached before attempting request
		select {
		case <-deadlineCtx.Done():
			log.Printf("[HTTP Retry] Deadline reached before attempt %d/%d for %s %s",
				attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String())
			if lastErr != nil {
				log.Printf("[HTTP Retry] Last error: %v", lastErr)
				return lastResp, lastErr
			}
			log.Printf("[HTTP Retry] Deadline context error: %v", deadlineCtx.Err())
			return nil, deadlineCtx.Err()
		default:
		}

		// Log attempt start
		if attempt > 0 {
			log.Printf("[HTTP Retry] Attempt %d/%d for %s %s",
				attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String())
		}

		// Create timeout context for this individual request attempt
		timeoutCtx, timeoutCancel := context.WithTimeout(deadlineCtx, timeout)

		// Update request context to use timeout context
		req = req.Clone(timeoutCtx)

		// Execute the request
		resp, err := http.DefaultClient.Do(req)

		// If successful, buffer the body before canceling timeout context
		if err == nil && resp != nil && resp.Body != nil {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()

			if readErr != nil {
				timeoutCancel()
				return nil, fmt.Errorf("[HTTP Retry] Failed to read response body: %w", readErr)
			}

			// Replace body with buffered version
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		timeoutCancel()

		// Check if we should retry
		if !shouldRetry(err, resp) {
			if attempt > 0 {
				log.Printf("[HTTP Retry] Request succeeded on attempt %d/%d for %s %s",
					attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String())
			}
			if err != nil {
				log.Printf("[HTTP Retry] Request failed on attempt %d/%d for %s %s: %v",
					attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String(), err)
				return nil, err
			}
			return resp, nil
		}

		// Log retry reason
		retryReason := getRetryReason(err, resp)
		log.Printf("[HTTP Retry] Retrying (attempt %d/%d) for %s %s: %s",
			attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String(), retryReason)

		// Store the error/response for potential return
		lastErr = err

		// If we got a response that needs retry, close it before retrying
		// (unless this is the last attempt, then we'll return it)
		if resp != nil && attempt < c.config.Retry.MaxRetries {
			// Close the response body to free resources before retry
			resp.Body.Close()
			log.Printf("[HTTP Retry] Closed response body for attempt %d/%d for %s %s",
				attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String())
		}

		if resp != nil {
			lastResp = resp
		}

		// If this was the last attempt, break
		if attempt == c.config.Retry.MaxRetries {
			log.Printf("[HTTP Retry] Max retries (%d) exhausted for %s %s",
				c.config.Retry.MaxRetries+1, req.Method, req.URL.String())
			break
		}

		// Wait before retrying (exponential backoff)
		// Check deadline context during backoff
		log.Printf("[HTTP Retry] Waiting %v before next retry for %s %s",
			backoff, req.Method, req.URL.String())
		select {
		case <-deadlineCtx.Done():
			log.Printf("[HTTP Retry] Deadline reached during backoff for %s %s",
				req.Method, req.URL.String())
			if lastErr != nil {
				return lastResp, lastErr
			}
			log.Printf("[HTTP Retry] Deadline reached during backoff for %s %s: %v",
				req.Method, req.URL.String(), deadlineCtx.Err())
			return nil, deadlineCtx.Err()
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
					log.Printf("[HTTP Retry] Failed to recreate request body for attempt %d/%d for %s %s: %v",
						attempt+1, c.config.Retry.MaxRetries+1, req.Method, req.URL.String(), err)
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

// getRetryReason returns a human-readable reason for why a retry is happening
func getRetryReason(err error, resp *http.Response) string {
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if urlErr.Timeout() {
				return "timeout error"
			}
			if urlErr.Temporary() {
				return "temporary error"
			}
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return "context deadline exceeded"
		}
		return err.Error()
	}

	if resp != nil && resp.StatusCode >= 500 {
		return "HTTP " + http.StatusText(resp.StatusCode)
	}

	return "unknown reason"
}
