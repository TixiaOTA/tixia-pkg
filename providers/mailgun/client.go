package mailgun

import (
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
)

// DefaultBaseURL is the Mailgun US API base URL.
const DefaultBaseURL = "https://api.mailgun.net"

// Config represents the Mailgun client configuration.
type Config struct {
	APIKey      string              // Mailgun private API key (required)
	Domain      string              // Mailgun sending domain (required)
	BaseURL     string              // Mailgun API base URL (optional, default: https://api.mailgun.net)
	Timeout     *time.Duration      // Timeout for individual request attempts (optional, default: 30s)
	RetryConfig *client.RetryConfig // Retry configuration for the HTTP client (optional, default: 2 retries)
}

// Validate checks that all required configuration fields are set.
func (c Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.Domain == "" {
		return ErrMissingDomain
	}

	return nil
}

// Client is a Mailgun API client.
type Client struct {
	config     Config
	httpClient *client.HTTPClient
	logger     *logger.Logger
}

// New creates a new Mailgun client.
func New(config Config, logger *logger.Logger) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")

	if config.RetryConfig == nil {
		config.RetryConfig = &client.RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    time.Second,
			MaxBackoff:        8 * time.Second,
			BackoffMultiplier: 2.0,
		}
	}

	httpClient := client.NewHTTPClient(client.HTTPConfig{
		BaseURL: config.BaseURL,
		Timeout: config.Timeout,
		Retry:   config.RetryConfig,
	})

	return &Client{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}
