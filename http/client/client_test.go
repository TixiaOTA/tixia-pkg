package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HTTPClientTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *HTTPClient
}

func (suite *HTTPClientTestSuite) SetupTest() {
	// Server will be created per test as needed
	suite.server = nil
	suite.client = nil
}

func (suite *HTTPClientTestSuite) TearDownTest() {
	if suite.server != nil {
		suite.server.Close()
	}
}

func TestHTTPClientTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientTestSuite))
}

// Helper function to create a test server
func (suite *HTTPClientTestSuite) createTestServer(handler http.HandlerFunc) {
	suite.server = httptest.NewServer(handler)
}

// Helper function to create a client with test server URL
func (suite *HTTPClientTestSuite) createClient(config HTTPConfig) {
	if suite.server != nil {
		config.BaseURL = suite.server.URL
	}
	suite.client = NewHTTPClient(config)
}

// =============================================================================
// Constructor & Configuration Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestNewHTTPClient_WithDefaultRetryConfig() {
	client := NewHTTPClient(HTTPConfig{
		BaseURL: "https://example.com",
	})

	suite.NotNil(client)
	suite.NotNil(client.config.Retry)
	suite.Equal(3, client.config.Retry.MaxRetries)
	suite.Equal(100*time.Millisecond, client.config.Retry.InitialBackoff)
	suite.Equal(5*time.Second, client.config.Retry.MaxBackoff)
	suite.Equal(2.0, client.config.Retry.BackoffMultiplier)
}

func (suite *HTTPClientTestSuite) TestNewHTTPClient_WithCustomRetryConfig() {
	customRetry := &RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    200 * time.Millisecond,
		MaxBackoff:        10 * time.Second,
		BackoffMultiplier: 3.0,
	}

	client := NewHTTPClient(HTTPConfig{
		BaseURL: "https://example.com",
		Retry:   customRetry,
	})

	suite.NotNil(client)
	suite.Equal(5, client.config.Retry.MaxRetries)
	suite.Equal(200*time.Millisecond, client.config.Retry.InitialBackoff)
	suite.Equal(10*time.Second, client.config.Retry.MaxBackoff)
	suite.Equal(3.0, client.config.Retry.BackoffMultiplier)
}

func (suite *HTTPClientTestSuite) TestSetBaseURL() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	client.SetBaseURL("https://newurl.com")
	suite.Equal("https://newurl.com", client.config.BaseURL)
}

func (suite *HTTPClientTestSuite) TestSetTimeout() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	timeout := 10 * time.Second
	client.SetTimeout(timeout)
	suite.NotNil(client.config.Timeout)
	suite.Equal(timeout, *client.config.Timeout)
}

func (suite *HTTPClientTestSuite) TestSetDeadline() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	deadline := 20 * time.Second
	client.SetDeadline(deadline)
	suite.NotNil(client.config.Deadline)
	suite.Equal(deadline, *client.config.Deadline)
}

func (suite *HTTPClientTestSuite) TestSetRetry() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	retry := RetryConfig{
		MaxRetries:        10,
		InitialBackoff:    500 * time.Millisecond,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 4.0,
	}
	client.SetRetry(retry)
	suite.Equal(10, client.config.Retry.MaxRetries)
	suite.Equal(500*time.Millisecond, client.config.Retry.InitialBackoff)
	suite.Equal(30*time.Second, client.config.Retry.MaxBackoff)
	suite.Equal(4.0, client.config.Retry.BackoffMultiplier)
}

func (suite *HTTPClientTestSuite) TestSetMaxRetries() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	client.SetMaxRetries(7)
	suite.Equal(7, client.config.Retry.MaxRetries)
}

func (suite *HTTPClientTestSuite) TestSetInitialBackoff() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	backoff := 300 * time.Millisecond
	client.SetInitialBackoff(backoff)
	suite.Equal(backoff, client.config.Retry.InitialBackoff)
}

func (suite *HTTPClientTestSuite) TestSetMaxBackoff() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	maxBackoff := 15 * time.Second
	client.SetMaxBackoff(maxBackoff)
	suite.Equal(maxBackoff, client.config.Retry.MaxBackoff)
}

func (suite *HTTPClientTestSuite) TestSetBackoffMultiplier() {
	client := NewHTTPClient(HTTPConfig{BaseURL: "https://example.com"})
	multiplier := 2.5
	client.SetBackoffMultiplier(multiplier)
	suite.Equal(multiplier, client.config.Retry.BackoffMultiplier)
}

// =============================================================================
// GET Request Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestGet_Success() {
	expectedBody := []byte(`{"message": "success"}`)
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodGet, r.Method)
		suite.Equal("/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write(expectedBody)
	})

	suite.createClient(HTTPConfig{})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)
	suite.Equal(expectedBody, body)
}

func (suite *HTTPClientTestSuite) TestGet_WithCustomHeaders() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("application/json", r.Header.Get("Accept"))
		suite.Equal("Bearer token123", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	suite.createClient(HTTPConfig{})

	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer token123",
	}

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path:    "/test",
		Headers: &headers,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestGet_WithTimeout() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "slow response"}`))
	})

	timeout := 50 * time.Millisecond
	suite.createClient(HTTPConfig{
		Timeout: &timeout,
	})

	ctx := context.Background()
	_, err := suite.client.Get(ctx, GetRequest{
		Path: "/slow",
	})

	suite.Error(err)
}

func (suite *HTTPClientTestSuite) TestGet_WithDeadline() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
	})

	timeout := 50 * time.Millisecond
	deadline := 300 * time.Millisecond
	suite.createClient(HTTPConfig{
		Timeout:  &timeout,
		Deadline: &deadline,
		Retry: &RetryConfig{
			MaxRetries:        10,
			InitialBackoff:    50 * time.Millisecond,
			MaxBackoff:        1 * time.Second,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	_, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.Error(err)
	suite.True(attemptCount < 10, "Should stop retrying when deadline is reached")
}

func (suite *HTTPClientTestSuite) TestGet_RetryOn5xxError() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success after retries"}`))
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(3, attemptCount)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestGet_MaxRetriesExhausted() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusInternalServerError)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	suite.Equal(3, attemptCount) // Initial + 2 retries
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestGet_NoRetryOn4xxError() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	suite.Equal(1, attemptCount) // Should not retry on 4xx
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestGet_ContextCancellation() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.Error(err)
	suite.Equal(context.Canceled, err)
}

// =============================================================================
// POST JSON Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestPostJSON_Success() {
	requestBody := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/users", r.URL.Path)
		suite.Equal("application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		var receivedBody map[string]interface{}
		json.Unmarshal(body, &receivedBody)
		suite.Equal("John Doe", receivedBody["name"])
		suite.Equal("john@example.com", receivedBody["email"])

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123", "status": "created"}`))
	})

	suite.createClient(HTTPConfig{})

	ctx := context.Background()
	resp, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path: "/users",
		Body: requestBody,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostJSON_WithCustomHeaders() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("Bearer token456", r.Header.Get("Authorization"))
		suite.Equal("application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{})

	headers := map[string]string{
		"Authorization": "Bearer token456",
	}

	ctx := context.Background()
	resp, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path:    "/api/data",
		Body:    map[string]string{"key": "value"},
		Headers: &headers,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostJSON_RetryOn5xxError() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path: "/api/data",
		Body: map[string]string{"test": "data"},
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(2, attemptCount)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostJSON_InvalidBody() {
	suite.createClient(HTTPConfig{
		BaseURL: "https://example.com",
	})

	// Channel cannot be marshaled to JSON
	invalidBody := make(chan int)

	ctx := context.Background()
	resp, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path: "/test",
		Body: invalidBody,
	})

	suite.Error(err)
	suite.Nil(resp)
}

func (suite *HTTPClientTestSuite) TestPostJSON_WithTimeout() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	timeout := 50 * time.Millisecond
	suite.createClient(HTTPConfig{
		Timeout: &timeout,
	})

	ctx := context.Background()
	_, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path: "/slow",
		Body: map[string]string{"test": "data"},
	})

	suite.Error(err)
}

// =============================================================================
// POST Form Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestPostForm_Success() {
	formData := url.Values{}
	formData.Set("username", "johndoe")
	formData.Set("password", "secret123")

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/login", r.URL.Path)
		suite.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := r.ParseForm()
		suite.NoError(err)
		suite.Equal("johndoe", r.FormValue("username"))
		suite.Equal("secret123", r.FormValue("password"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token": "abc123"}`))
	})

	suite.createClient(HTTPConfig{})

	ctx := context.Background()
	resp, err := suite.client.PostForm(ctx, PostFormRequest{
		Path: "/login",
		Body: formData,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostForm_WithCustomHeaders() {
	formData := url.Values{}
	formData.Set("field", "value")

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("custom-value", r.Header.Get("X-Custom-Header"))
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{})

	headers := map[string]string{
		"X-Custom-Header": "custom-value",
	}

	ctx := context.Background()
	resp, err := suite.client.PostForm(ctx, PostFormRequest{
		Path:    "/form",
		Body:    formData,
		Headers: &headers,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostForm_RetryOn5xxError() {
	attemptCount := 0
	formData := url.Values{}
	formData.Set("data", "test")

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.PostForm(ctx, PostFormRequest{
		Path: "/form",
		Body: formData,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(3, attemptCount)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostForm_WithTimeout() {
	formData := url.Values{}
	formData.Set("field", "value")

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	timeout := 50 * time.Millisecond
	suite.createClient(HTTPConfig{
		Timeout: &timeout,
	})

	ctx := context.Background()
	_, err := suite.client.PostForm(ctx, PostFormRequest{
		Path: "/slow",
		Body: formData,
	})

	suite.Error(err)
}

// =============================================================================
// Retry Logic & Backoff Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestRetry_ExponentialBackoff() {
	attemptCount := 0
	attemptTimes := []time.Time{}

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		attemptTimes = append(attemptTimes, time.Now())
		if attemptCount < 4 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    50 * time.Millisecond,
			MaxBackoff:        500 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(4, attemptCount)

	// Verify exponential backoff timing
	if len(attemptTimes) >= 2 {
		firstBackoff := attemptTimes[1].Sub(attemptTimes[0])
		suite.GreaterOrEqual(firstBackoff.Milliseconds(), int64(50))
	}

	if len(attemptTimes) >= 3 {
		secondBackoff := attemptTimes[2].Sub(attemptTimes[1])
		suite.GreaterOrEqual(secondBackoff.Milliseconds(), int64(100))
	}
}

func (suite *HTTPClientTestSuite) TestRetry_MaxBackoffLimit() {
	attemptCount := 0
	attemptTimes := []time.Time{}

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		attemptTimes = append(attemptTimes, time.Now())
		if attemptCount < 5 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        5,
			InitialBackoff:    50 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond, // Max backoff is low
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(5, attemptCount)

	// After initial backoff (50ms) and second backoff (100ms),
	// all subsequent backoffs should be capped at MaxBackoff (100ms)
	if len(attemptTimes) >= 4 {
		thirdBackoff := attemptTimes[3].Sub(attemptTimes[2])
		// Should be capped at 100ms, not 200ms
		suite.LessOrEqual(thirdBackoff.Milliseconds(), int64(150)) // Allow some margin
	}
}

func (suite *HTTPClientTestSuite) TestRetry_DeadlineDuringBackoff() {
	attemptCount := 0
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusInternalServerError)
	})

	deadline := 150 * time.Millisecond
	suite.createClient(HTTPConfig{
		Deadline: &deadline,
		Retry: &RetryConfig{
			MaxRetries:        10,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        1 * time.Second,
			BackoffMultiplier: 2.0,
		},
	})

	ctx := context.Background()
	_, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	// When deadline is reached during backoff, context.DeadlineExceeded is returned
	suite.Error(err)
	suite.ErrorIs(err, context.DeadlineExceeded)
	suite.True(attemptCount < 10, "Should stop retrying when deadline is reached during backoff")
}

// =============================================================================
// Edge Cases & Error Handling Tests
// =============================================================================

func (suite *HTTPClientTestSuite) TestGet_InvalidURL() {
	suite.createClient(HTTPConfig{
		BaseURL: "://invalid-url",
	})

	ctx := context.Background()
	_, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.Error(err)
}

func (suite *HTTPClientTestSuite) TestGet_EmptyResponse() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	suite.createClient(HTTPConfig{})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/empty",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusNoContent, resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)
	suite.Empty(body)
}

func (suite *HTTPClientTestSuite) TestPostJSON_BodyRereadingOnRetry() {
	attemptCount := 0
	receivedBodies := []string{}

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		body, _ := io.ReadAll(r.Body)
		receivedBodies = append(receivedBodies, string(body))

		if attemptCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	requestBody := map[string]string{"test": "data"}
	ctx := context.Background()
	resp, err := suite.client.PostJSON(ctx, PostJSONRequest{
		Path: "/test",
		Body: requestBody,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(2, attemptCount)

	// Verify body was re-read correctly on retry
	suite.Equal(receivedBodies[0], receivedBodies[1])
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestPostForm_BodyRereadingOnRetry() {
	attemptCount := 0
	receivedBodies := []string{}

	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		body, _ := io.ReadAll(r.Body)
		receivedBodies = append(receivedBodies, string(body))

		if attemptCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	suite.createClient(HTTPConfig{
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	})

	formData := url.Values{}
	formData.Set("field1", "value1")
	formData.Set("field2", "value2")

	ctx := context.Background()
	resp, err := suite.client.PostForm(ctx, PostFormRequest{
		Path: "/test",
		Body: formData,
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(2, attemptCount)

	// Verify body was re-read correctly on retry
	suite.Equal(receivedBodies[0], receivedBodies[1])
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestGet_DefaultTimeout() {
	suite.createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Create client without explicit timeout
	suite.createClient(HTTPConfig{})

	ctx := context.Background()
	resp, err := suite.client.Get(ctx, GetRequest{
		Path: "/test",
	})

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *HTTPClientTestSuite) TestRetry_Different5xxStatusCodes() {
	testCases := []struct {
		name        string
		statusCode  int
		shouldRetry bool
	}{
		{"InternalServerError", http.StatusInternalServerError, true},
		{"BadGateway", http.StatusBadGateway, true},
		{"ServiceUnavailable", http.StatusServiceUnavailable, true},
		{"GatewayTimeout", http.StatusGatewayTimeout, true},
		{"BadRequest", http.StatusBadRequest, false},
		{"NotFound", http.StatusNotFound, false},
		{"OK", http.StatusOK, false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			attemptCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attemptCount++
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			client := NewHTTPClient(HTTPConfig{
				BaseURL: server.URL,
				Retry: &RetryConfig{
					MaxRetries:        2,
					InitialBackoff:    10 * time.Millisecond,
					MaxBackoff:        100 * time.Millisecond,
					BackoffMultiplier: 2.0,
				},
			})

			ctx := context.Background()
			resp, err := client.Get(ctx, GetRequest{
				Path: "/test",
			})

			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tc.statusCode, resp.StatusCode)
			resp.Body.Close()

			if tc.shouldRetry {
				suite.Equal(3, attemptCount, fmt.Sprintf("%s should trigger retries", tc.name))
			} else {
				suite.Equal(1, attemptCount, fmt.Sprintf("%s should not trigger retries", tc.name))
			}
		})
	}
}
