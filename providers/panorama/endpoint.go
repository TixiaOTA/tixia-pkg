package panorama

import (
	"context"
	"io"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

const (
	sessionPath        = "/session"
	bookingAirlinePath = "/booking/airline"
)

const (
	signInTimeout     = 10 * time.Second
	getBookingTimeout = 10 * time.Second
	cancelTimeout     = 10 * time.Second
	searchTimeout     = 90 * time.Second
)

// SignIn signs in to the Panorama API
func (c *Client) SignIn(ctx context.Context) (*SignInResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "SignIn")
	defer span.End()

	// Convert the request to the internal model
	internalReq := c.config.Credentials // Sign In Request is the same as the credentials

	// Sanitize the request for logging
	sanitizedReq, _ := logger.Sanitize(internalReq)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SignIn request", zap.Any("request", sanitizedReq))

	// Set timeout
	c.httpClient.SetTimeout(signInTimeout)

	// Add header
	headers := map[string]string{
		"Action": string(SignIn),
	}

	httpResp, err := c.sendRequest(ctx, internalReq, sessionPath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp SignInResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// GetBooking gets a booking by order number
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "Search")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting Search request", zap.Any("request", req))

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	// Set timeout
	c.httpClient.SetTimeout(searchTimeout)

	// Add header
	headers := map[string]string{
		"Action":    string(Search),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers)
	if err != nil {
		// api error

		// panorma err
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp SearchResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// GetBooking gets a booking by order number
func (c *Client) GetBooking(ctx context.Context, req GetBookingRequest) (*GetBookingResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetBooking")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetBooking request", zap.Any("request", req))

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	// Set timeout
	c.httpClient.SetTimeout(getBookingTimeout)

	// Add header
	headers := map[string]string{
		"Action":    string(GetBooking),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp GetBookingResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// CancelVoid cancels a booking by order number with a booking state of "booked" to state of "void"
func (c *Client) CancelVoid(ctx context.Context, req CancelVoidRequest) (*GetBookingResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "CancelVoid")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CancelVoid request", zap.Any("request", req))

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	// Set timeout
	c.httpClient.SetTimeout(cancelTimeout)

	// Add header
	headers := map[string]string{
		"Action":    string(Cancel),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp GetBookingResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}
