package xendit

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
)

const defaultBaseURL = "https://api.xendit.co"

type Config struct {
	BaseURL      string              // Base URL for the API (optional, default: https://api.xendit.co)
	SecretAPIKey string              // Secret API key for Basic Auth (required)
	ForUserID    string              // XenPlatform sub-account user id (optional)
	SplitRuleID  string              // XenPlatform split rule id (optional)
	RetryConfig  *client.RetryConfig // Retry configuration (optional)
}

type Client struct {
	config         Config
	httpClient     *client.HTTPClient
	logger         *logger.Logger
	authHeader     string
}

func New(config Config, logger *logger.Logger) (*Client, error) {
	if config.SecretAPIKey == "" {
		return nil, errors.New("secret API key is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

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

	authValue := base64.StdEncoding.EncodeToString([]byte(config.SecretAPIKey + ":"))

	return &Client{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
		authHeader: "Basic " + authValue,
	}, nil
}
