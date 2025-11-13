package mg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	httppkg "github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

// MGClient represents a client for interacting with MG hotel API
type MGClient struct {
	config     MGConfig
	logger     *logger.Logger
	httpClient *httppkg.HTTPClient
}

// NewMGClient creates a new MG client instance
// Returns an error if the logger is nil or if the config validation fails
func NewMGClient(config MGConfig, logger *logger.Logger) (*MGClient, error) {
	// Check if logger is nil
	if logger == nil {
		return nil, fmt.Errorf("logger is required for MG client")
	}

	// Validate config
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate MG config: %w", err)
	}

	// Set default timeout if not provided
	timeout := 10 * time.Second
	if config.Timeout != nil {
		timeout = *config.Timeout
	}

	// Set default retry config if not provided
	retryConfig := httppkg.RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    500 * time.Millisecond,
		MaxBackoff:        5 * time.Second,
		BackoffMultiplier: 2.0,
	}
	if config.Retry != nil {
		retryConfig = *config.Retry
	}

	httpClient := httppkg.NewHTTPClient(httppkg.HTTPConfig{
		BaseURL: config.BaseURL,
		Timeout: &timeout,
		Retry:   &retryConfig,
	})

	return &MGClient{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}, nil
}

// SearchHotelAvailability performs hotel availability search
func (c *MGClient) SearchHotelAvailability(ctx context.Context, request *HotelSearchRequest) (HotelSearchResponse, error) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "SearchHotelAvailability", "SearchHotelAvailability")
	defer span.End()

	var response HotelSearchResponse

	// Validate request
	if err := request.Validate(); err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("failed to validate request", zap.Error(err))
		return response, fmt.Errorf("failed to validate request: %w", err)
	}

	// Set login credentials from config
	if request.Login == nil {
		request.Login = &Login{
			Username:   c.config.LoginName,
			Password:   c.config.Password,
			AgencyCode: c.config.AgentID,
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("Sending hotel search request to MG",
		zap.Any("request", request),
	)

	resp, err := c.httpClient.PostJSON(ctx, httppkg.PostJSONRequest{
		Path:    SearchHotelURL,
		Body:    request,
		Headers: &DefaultHeaders,
	})
	if err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("failed to send hotel search request", zap.Error(err))
		return response, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("failed to read response body", zap.Error(err))
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal response body to HotelSearchResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Warn("response body is not JSON", zap.Any("response", response), zap.Error(err))
		return response, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// Check if status is false and error code is JRVXML161, if so, retry the request
	if !response.Status {
		if response.ErrorCode != nil && *response.ErrorCode == "JRVXML161" {
			c.logger.WithTrace(ctx, currentSpanId).Warn("provider returned JRVXML161, retrying", zap.Any("response", response))
			return c.SearchHotelAvailability(ctx, request)
		}
	}

	return response, nil
}
