package mailgun

import (
	"testing"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name    string
		config  Config
		wantErr error
	}{
		{
			name:    "missing api key",
			config:  Config{Domain: "example.com"},
			wantErr: ErrMissingAPIKey,
		},
		{
			name:    "missing domain",
			config:  Config{APIKey: "key"},
			wantErr: ErrMissingDomain,
		},
		{
			name:   "valid config",
			config: Config{APIKey: "key", Domain: "example.com"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := New(testCase.config, logger.Default())

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestNew_NilLogger(t *testing.T) {
	_, err := New(Config{APIKey: "key", Domain: "example.com"}, nil)
	require.ErrorIs(t, err, ErrNilLogger)
}

func TestNew_DefaultBaseURLAndRetry(t *testing.T) {
	got, err := New(Config{APIKey: "key", Domain: "example.com"}, logger.Default())
	require.NoError(t, err)

	assert.Equal(t, DefaultBaseURL, got.config.BaseURL)
	require.NotNil(t, got.config.RetryConfig)
	assert.Equal(t, 2, got.config.RetryConfig.MaxRetries)
	assert.Equal(t, time.Second, got.config.RetryConfig.InitialBackoff)
	assert.Equal(t, 8*time.Second, got.config.RetryConfig.MaxBackoff)
	assert.Equal(t, 2.0, got.config.RetryConfig.BackoffMultiplier)
}

func TestNew_TrimsTrailingSlashFromBaseURL(t *testing.T) {
	got, err := New(Config{
		APIKey:  "key",
		Domain:  "example.com",
		BaseURL: "https://api.eu.mailgun.net/",
	}, logger.Default())
	require.NoError(t, err)

	assert.Equal(t, "https://api.eu.mailgun.net", got.config.BaseURL)
}

func TestNew_CustomRetryConfig(t *testing.T) {
	retryConfig := &client.RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        time.Second,
		BackoffMultiplier: 3.0,
	}

	got, err := New(Config{
		APIKey:      "key",
		Domain:      "example.com",
		RetryConfig: retryConfig,
	}, logger.Default())
	require.NoError(t, err)

	assert.Same(t, retryConfig, got.config.RetryConfig)
}

func TestSendEmailRequest_Validate(t *testing.T) {
	validRequest := SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Subject",
		Text:    "Hello",
	}

	testCases := []struct {
		name    string
		request SendEmailRequest
		wantErr error
	}{
		{
			name:    "valid request",
			request: validRequest,
		},
		{
			name:    "missing from",
			request: SendEmailRequest{To: validRequest.To, Subject: "Subject", Text: "Hello"},
			wantErr: ErrMissingFrom,
		},
		{
			name:    "missing recipients",
			request: SendEmailRequest{From: validRequest.From, Subject: "Subject", Text: "Hello"},
			wantErr: ErrMissingRecipients,
		},
		{
			name:    "missing subject",
			request: SendEmailRequest{From: validRequest.From, To: validRequest.To, Text: "Hello"},
			wantErr: ErrMissingSubject,
		},
		{
			name:    "missing body",
			request: SendEmailRequest{From: validRequest.From, To: validRequest.To, Subject: "Subject"},
			wantErr: ErrMissingBody,
		},
		{
			name: "invalid attachment",
			request: SendEmailRequest{
				From:    validRequest.From,
				To:      validRequest.To,
				Subject: "Subject",
				Text:    "Hello",
				Attachments: []Attachment{
					{Filename: "", Content: []byte("data")},
				},
			},
			wantErr: ErrInvalidAttachment,
		},
		{
			name: "valid html request with attachment",
			request: SendEmailRequest{
				From:    validRequest.From,
				To:      validRequest.To,
				Subject: "Subject",
				HTML:    "<p>Hello</p>",
				Attachments: []Attachment{
					{Filename: "file.txt", Content: []byte("data")},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.request.Validate()

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
