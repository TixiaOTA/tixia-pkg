package vervotech

import (
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

const (
	DefaultBaseURL    = "https://hotelmapping.vervotech.com/api"
	DefaultAPIVersion = "3.0"
)

type Config struct {
	BaseURL    string // Base URL for the API (default: "https://hotelmapping.vervotech.com/api")
	APIKey     string // API key for authentication (required)
	AccountID  string // Account ID used in request headers (required)
	APIVersion string // API version used in URL path (default: "3.0")
	SyncID     string // Sync ID used in some request headers (optional)

	RetryConfig *client.RetryConfig // Retry configuration for the HTTP client (optional, default: 3 retries, 1 second initial backoff, 5 second max backoff, 2x multiplier)
}

type Client struct {
	config     Config
	httpClient *client.HTTPClient
}

func New(config Config) (*Client, error) {
	// Set defaults if not provided
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.APIVersion == "" {
		config.APIVersion = DefaultAPIVersion
	}

	// Set default HTTP config if not provided
	if config.RetryConfig == nil {
		config.RetryConfig = &client.RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    1 * time.Second,
			MaxBackoff:        5 * time.Second,
			BackoffMultiplier: 2.0,
		}
	}

	httpClient := client.NewHTTPClient(client.HTTPConfig{
		BaseURL: config.BaseURL,
		Retry:   config.RetryConfig,
	})

	return &Client{
		config:     config,
		httpClient: httpClient,
	}, nil
}

func (c *Client) SetSyncID(syncID string) {
	c.config.SyncID = syncID
}
