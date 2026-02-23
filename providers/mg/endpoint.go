package mg

import (
	"context"
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
	getHotelDetailPath    = "/GetHotelDetail"
	getDestinationsPath   = "/GetDestinations"
)

const (
	searchHotelTimeout       = 10 * time.Second
	recheckHotelTimeout      = 10 * time.Second
	bookHotelTimeout         = 120 * time.Second
	getRSVNDetailsTimeout    = 10 * time.Second
	cancelReservationTimeout = 10 * time.Second
	getRSVNListTimeout       = 10 * time.Second
	getHotelDetailTimeout    = 10 * time.Second
	getDestinationsTimeout   = 10 * time.Second
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
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp SearchHotelResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

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
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp RecheckHotelResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

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
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp BookHotelResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

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
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp GetRSVNDetailsResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

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
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp CancelReservationResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

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

	if err := req.Validate(); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    BadRequestError,
			ErrorMessage: err.Error(),
		}
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

	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp GetRSVNListResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

// GetHotelDetail gets the hotel details
func (c *Client) GetHotelDetail(ctx context.Context, req GetHotelDetailRequestDTO) (*GetHotelDetailResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "GetHotelDetail")
	defer span.End()

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    BadRequestError,
			ErrorMessage: err.Error(),
		}
	}

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetHotelDetail request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(getHotelDetailTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, getHotelDetailPath, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp GetHotelDetailResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}

// GetDestinations gets the destinations
func (c *Client) GetDestinations(ctx context.Context, req GetDestinationsRequestDTO) (*GetDestinationsResponse, *MGError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "mg", "GetDestinations")
	defer span.End()

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    BadRequestError,
			ErrorMessage: err.Error(),
		}
	}

	// Convert the request to the internal model
	internalReq := req.ToInternalModel(c.config.Credentials)

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetDestinations request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(getDestinationsTimeout)

	httpResp, err := c.sendRequest(ctx, internalReq, getDestinationsPath, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	var resp GetDestinationsResponse
	if err := parseResponseBody(httpResp.Body, &resp); err != nil {
		if _, ok := err.(*MGError); ok {
			return nil, err.(*MGError)
		}

		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}
