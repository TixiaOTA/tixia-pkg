package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func TestPostJSON_RequestTimeoutOverride(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	clientTimeout := time.Second
	requestTimeout := 10 * time.Millisecond
	client := NewHTTPClient(HTTPConfig{
		BaseURL: server.URL,
		Timeout: &clientTimeout,
		Retry:   &RetryConfig{MaxRetries: 0},
	})

	_, err := client.PostJSON(context.Background(), PostJSONRequest{
		Path: "/slow",
		Body: map[string]string{"test": "data"},
		Options: &RequestOptions{
			Timeout: &requestTimeout,
		},
	})
	if err == nil {
		t.Fatal("expected request timeout error")
	}
}

func TestPostJSON_RequestRetryOverride(t *testing.T) {
	attempts := 0
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	client := NewHTTPClient(HTTPConfig{
		BaseURL: server.URL,
		Retry:   &RetryConfig{MaxRetries: 0},
	})
	requestRetry := &RetryConfig{
		MaxRetries:        1,
		InitialBackoff:    time.Millisecond,
		MaxBackoff:        time.Millisecond,
		BackoffMultiplier: 1,
	}

	resp, err := client.PostJSON(context.Background(), PostJSONRequest{
		Path: "/retry",
		Body: map[string]string{"test": "data"},
		Options: &RequestOptions{
			Retry: requestRetry,
		},
	})
	if err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	defer resp.Body.Close()

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if client.config.Retry.MaxRetries != 0 {
		t.Fatalf("request retry override mutated client config: %d", client.config.Retry.MaxRetries)
	}
}

func TestGet_NoRetriesBuffersStreamedResponse(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("first "))
		w.(http.Flusher).Flush()
		time.Sleep(50 * time.Millisecond)
		w.Write([]byte("second"))
	})
	defer server.Close()

	timeout := time.Second
	client := NewHTTPClient(HTTPConfig{
		BaseURL: server.URL,
		Timeout: &timeout,
		Retry:   &RetryConfig{MaxRetries: 0},
	})

	resp, err := client.Get(context.Background(), GetRequest{Path: "/stream"})
	if err != nil {
		t.Fatalf("Get streamed response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read buffered body: %v", err)
	}
	if string(body) != "first second" {
		t.Fatalf("body = %q, want %q", body, "first second")
	}
}

func TestGet_NoRetriesPreservesBodyReadError(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "10")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("short"))
	})
	defer server.Close()

	client := NewHTTPClient(HTTPConfig{
		BaseURL: server.URL,
		Retry:   &RetryConfig{MaxRetries: 0},
	})
	resp, err := client.Get(context.Background(), GetRequest{Path: "/short"})
	if resp != nil {
		resp.Body.Close()
		t.Fatalf("response = %v, want nil on body read error", resp)
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("error = %v, want wrapped unexpected EOF", err)
	}
}
