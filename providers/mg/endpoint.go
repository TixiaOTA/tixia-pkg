package mg

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

const (
	searchHotelPath  = "/SearchHotel"
	recheckHotelPath = "/RecheckHotel"
	bookHotelPath    = "/BookHotel"
)

const (
	searchHotelTimeout  = 10 * time.Second
	recheckHotelTimeout = 10 * time.Second
	bookHotelTimeout    = 120 * time.Second
)

// SearchHotel searches for hotels
func (c *Client) SearchHotel(ctx context.Context, req SearchHotelRequest) (*SearchHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "SearchHotel")
	defer span.End()

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SearchHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(searchHotelTimeout)

	httpResp, err := c.sendRequest(ctx, req, searchHotelPath, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	// Try to unmarshal as MGError first
	var mgError MGError
	if err := json.Unmarshal(bodyBytes, &mgError); err == nil {
		// Check if it's a valid MGError
		// Status = false, then it's an error
		if !mgError.Status || mgError.ErrorCode != "" || mgError.ErrorMessage != "" {
			return nil, &MGError{
				Status:       false,
				ErrorCode:    mgError.ErrorCode,
				ErrorMessage: mgError.ErrorMessage,
			}
		}
	}

	// If not an MGError, unmarshal as SearchHotelResponse
	var resp SearchHotelResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

// RecheckHotel rechecks a hotel
func (c *Client) RecheckHotel(ctx context.Context, req RecheckHotelRequest) (*RecheckHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "RecheckHotel")
	defer span.End()

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting RecheckHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(recheckHotelTimeout)

	httpResp, err := c.sendRequest(ctx, req, recheckHotelPath, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	// Try to unmarshal as MGError first
	var mgError MGError
	if err := json.Unmarshal(bodyBytes, &mgError); err == nil {
		// Check if it's a valid MGError
		// Status = false, then it's an error
		if !mgError.Status || mgError.ErrorCode != "" || mgError.ErrorMessage != "" {
			sanitizedReq, _ := logger.Sanitize(req)
			c.logger.Error("RecheckHotel failed", zap.Any("request", sanitizedReq))
			return nil, &MGError{
				Status:       false,
				ErrorCode:    mgError.ErrorCode,
				ErrorMessage: mgError.ErrorMessage,
			}
		}
	}

	// If not an MGError, unmarshal as RecheckHotelResponse
	var resp RecheckHotelResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

// BookHotel books a hotel
func (c *Client) BookHotel(ctx context.Context, req BookHotelRequest) (*BookHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "BookHotel")
	defer span.End()

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting BookHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(bookHotelTimeout)
	c.httpClient.SetMaxRetries(0) // No retries for BookHotel

	httpResp, err := c.sendRequest(ctx, req, bookHotelPath, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	// Try to unmarshal as MGError first
	var mgError MGError
	if err := json.Unmarshal(bodyBytes, &mgError); err == nil {
		// Check if it's a valid MGError
		// Status = false, then it's an error
		if !mgError.Status || mgError.ErrorCode != "" || mgError.ErrorMessage != "" {
			return nil, &MGError{
				Status:       false,
				ErrorCode:    mgError.ErrorCode,
				ErrorMessage: mgError.ErrorMessage,
			}
		}
	}

	// If not an MGError, unmarshal as BookHotelResponse
	var resp BookHotelResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}
