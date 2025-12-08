package panorama

import (
	"errors"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
)

type Config struct {
	BaseURL     string      // Base URL for the API (required)
	Credentials Credentials // Credentials for authentication (required)

	RetryConfig *client.RetryConfig // Retry configuration for the HTTP client (optional, default: 3 retries, 5 second initial backoff, 15 second max backoff, 2x multiplier)
}

type Client struct {
	config     Config
	httpClient *client.HTTPClient
	logger     *logger.Logger
}

func New(config Config, logger *logger.Logger) (*Client, error) {
	// Set defaults if not provided
	if config.BaseURL == "" {
		return nil, errors.New("base URL is required")
	}

	// Validate credentials
	if err := config.Credentials.Validate(); err != nil {
		return nil, err
	}

	// Set default HTTP config if not provided
	if config.RetryConfig == nil {
		config.RetryConfig = &client.RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    5 * time.Second,
			MaxBackoff:        15 * time.Second,
			BackoffMultiplier: 2.0,
		}
	}

	httpClient := client.NewHTTPClient(client.HTTPConfig{
		BaseURL: config.BaseURL,
		Retry:   config.RetryConfig,
	})

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}
