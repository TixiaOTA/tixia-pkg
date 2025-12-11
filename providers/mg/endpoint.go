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
	searchHotelPath       = "/SearchHotel"
	recheckHotelPath      = "/RecheckHotel"
	bookHotelPath         = "/BookHotel"
	getRSVNDetailsPath    = "/GetRSVNDetails"
	cancelReservationPath = "/CancelReservation"
	getRSVNListPath       = "/GetRSVNList"
)

const (
	searchHotelTimeout       = 10 * time.Second
	recheckHotelTimeout      = 10 * time.Second
	bookHotelTimeout         = 120 * time.Second
	getRSVNDetailsTimeout    = 10 * time.Second
	cancelReservationTimeout = 10 * time.Second
	getRSVNListTimeout       = 10 * time.Second
)

// SearchHotel searches for hotels
func (c *Client) SearchHotel(ctx context.Context, req SearchHotelRequestDTO) (*SearchHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "SearchHotel")
	defer span.End()

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SearchHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(searchHotelTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, searchHotelPath, nil)
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
func (c *Client) RecheckHotel(ctx context.Context, req RecheckHotelRequestDTO) (*RecheckHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "RecheckHotel")
	defer span.End()

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting RecheckHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(recheckHotelTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, recheckHotelPath, nil)
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
func (c *Client) BookHotel(ctx context.Context, req BookHotelRequestDTO) (*BookHotelResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "BookHotel")
	defer span.End()

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting BookHotel request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(bookHotelTimeout)
	c.httpClient.SetMaxRetries(0) // No retries for BookHotel

	httpResp, err := c.sendRequest(ctx, internalReq, bookHotelPath, nil)
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

// GetRSVNDetails gets the RSVN details
func (c *Client) GetRSVNDetails(ctx context.Context, req GetRSVNDetailsRequestDTO) (*GetRSVNDetailsResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "GetRSVNDetails")
	defer span.End()

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetRSVNDetails request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(getRSVNDetailsTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, getRSVNDetailsPath, nil)
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

	// If not an MGError, unmarshal as GetRSVNDetailsResponse
	var resp GetRSVNDetailsResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

// CancelReservation cancels a reservation
func (c *Client) CancelReservation(ctx context.Context, req CancelReservationRequestDTO) (*CancelReservationResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "CancelReservation")
	defer span.End()

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CancelReservation request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(cancelReservationTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, cancelReservationPath, nil)
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
				RawBody:      &bodyBytes,
			}
		}
	}

	// If not an MGError, unmarshal as CancelReservationResponse
	var resp CancelReservationResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

func (c *Client) GetRSVNList(ctx context.Context, req GetRSVNListRequestDTO) (*GetRSVNListResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "GetRSVNList")
	defer span.End()

	// set 1 to get by Booking date and 2 to get by Check-in date.
	if req.DateType == "" {
		req.DateType = "1"
	}

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetRSVNList request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(getRSVNListTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, getRSVNListPath, nil)
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

	// If not an MGError, unmarshal as GetRSVNDetailsResponse
	var resp GetRSVNListResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}
