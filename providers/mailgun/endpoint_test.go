package mailgun

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	mailgunClient, err := New(Config{
		APIKey:      "key-test",
		Domain:      "example.com",
		BaseURL:     baseURL,
		RetryConfig: &client.RetryConfig{MaxRetries: 0},
	}, logger.Default())
	require.NoError(t, err)

	return mailgunClient
}

func TestClient_SendEmail_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v3/example.com/messages", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "api", username)
		assert.Equal(t, "key-test", password)

		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "sender@example.com", r.PostForm.Get("from"))
		assert.Equal(t, []string{"one@example.com", "two@example.com"}, r.PostForm["to"])
		assert.Equal(t, []string{"cc@example.com"}, r.PostForm["cc"])
		assert.Equal(t, []string{"bcc@example.com"}, r.PostForm["bcc"])
		assert.Equal(t, "Hello", r.PostForm.Get("subject"))
		assert.Equal(t, "Plain body", r.PostForm.Get("text"))
		assert.Equal(t, "<p>HTML body</p>", r.PostForm.Get("html"))
		assert.Equal(t, "reply@example.com", r.PostForm.Get("h:Reply-To"))
		assert.Equal(t, "spring-sale", r.PostForm.Get("h:X-Campaign"))
		assert.Equal(t, []string{"welcome", "onboarding"}, r.PostForm["o:tag"])
		assert.Equal(t, `"Ada"`, r.PostForm.Get("v:name"))
		assert.Equal(t, "yes", r.PostForm.Get("o:testmode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendEmailResponse{
			ID:      "<message-id@example.com>",
			Message: "Queued. Thank you.",
		})
	}))
	defer server.Close()

	mailgunClient := newTestClient(t, server.URL)

	got, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"one@example.com", "two@example.com"},
		CC:      []string{"cc@example.com"},
		BCC:     []string{"bcc@example.com"},
		Subject: "Hello",
		Text:    "Plain body",
		HTML:    "<p>HTML body</p>",
		ReplyTo: "reply@example.com",
		Tags:    []string{"welcome", "onboarding"},
		Headers: map[string]string{"X-Campaign": "spring-sale"},
		Variables: map[string]any{
			"name": "Ada",
		},
		TestMode: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "<message-id@example.com>", got.ID)
	assert.Equal(t, "Queued. Thank you.", got.Message)
}

func TestClient_SendEmail_WithAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/example.com/messages", r.URL.Path)
		assert.True(t, strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary="))

		assert.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "sender@example.com", r.FormValue("from"))
		assert.Equal(t, []string{"recipient@example.com"}, r.MultipartForm.Value["to"])
		assert.Equal(t, "With attachment", r.FormValue("subject"))
		assert.Equal(t, "See attached", r.FormValue("text"))

		file, header, err := r.FormFile("attachment")
		assert.NoError(t, err)
		if assert.NotNil(t, file) {
			defer file.Close()

			assert.Equal(t, "receipt.txt", header.Filename)
			content, readErr := io.ReadAll(file)
			assert.NoError(t, readErr)
			assert.Equal(t, "attachment content", string(content))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendEmailResponse{Message: "Queued. Thank you."})
	}))
	defer server.Close()

	mailgunClient := newTestClient(t, server.URL)

	got, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "With attachment",
		Text:    "See attached",
		Attachments: []Attachment{
			{
				Filename:    "receipt.txt",
				Content:     []byte("attachment content"),
				ContentType: "text/plain",
			},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "Queued. Thank you.", got.Message)
}

func TestClient_SendEmail_InlineAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseMultipartForm(1<<20))

		file, header, err := r.FormFile("inline")
		assert.NoError(t, err)
		if assert.NotNil(t, file) {
			defer file.Close()
			assert.Equal(t, "logo.png", header.Filename)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendEmailResponse{Message: "Queued. Thank you."})
	}))
	defer server.Close()

	mailgunClient := newTestClient(t, server.URL)

	_, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Inline image",
		HTML:    `<img src="cid:logo.png">`,
		Inline: []Attachment{
			{Filename: "logo.png", Content: []byte("\x89PNG\r\n\x1a\n")},
		},
	})

	require.NoError(t, err)
}

func TestClient_SendEmail_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "'to' parameter is not a valid address",
		})
	}))
	defer server.Close()

	mailgunClient := newTestClient(t, server.URL)

	got, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"invalid"},
		Subject: "Hello",
		Text:    "Body",
	})

	require.Error(t, err)
	assert.Nil(t, got)

	var apiError *APIError
	require.True(t, errors.As(err, &apiError))
	assert.Equal(t, http.StatusBadRequest, apiError.StatusCode)
	assert.Equal(t, "'to' parameter is not a valid address", apiError.Message)
	assert.NotEmpty(t, apiError.RawBody)
}

func TestClient_SendEmail_NonJSONAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Forbidden"))
	}))
	defer server.Close()

	mailgunClient := newTestClient(t, server.URL)

	_, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Hello",
		Text:    "Body",
	})

	var apiError *APIError
	require.True(t, errors.As(err, &apiError))
	assert.Equal(t, http.StatusUnauthorized, apiError.StatusCode)
	assert.Equal(t, "Forbidden", apiError.Message)
}

func TestClient_SendEmail_ValidationError(t *testing.T) {
	mailgunClient := newTestClient(t, "http://127.0.0.1:1")

	got, err := mailgunClient.SendEmail(context.Background(), SendEmailRequest{})

	require.ErrorIs(t, err, ErrMissingFrom)
	assert.Nil(t, got)
}
