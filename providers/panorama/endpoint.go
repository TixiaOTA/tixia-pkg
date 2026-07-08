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
	accountPath        = "/account"
)

const (
	signInTimeout            = 10 * time.Second
	getBalanceTimeout        = 10 * time.Second
	getBookingTimeout        = 10 * time.Second
	cancelTimeout            = 10 * time.Second
	getPriceItineraryTimeout = 30 * time.Second
	sellJourneysTimeout      = 30 * time.Second
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

// GetBalance gets the account balance
func (c *Client) GetBalance(ctx context.Context, req GetBalanceRequest) (*GetBalanceResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetBalance")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetBalance request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(getBalanceTimeout)

	headers := map[string]string{
		"Action":    string(GetBalance),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, accountPath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp GetBalanceResponse
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

// GetPriceItinerary gets the price for a selected itinerary
func (c *Client) GetPriceItinerary(ctx context.Context, req GetPriceItineraryRequest) (*GetPriceItineraryResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetPriceItinerary")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetPriceItinerary request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(getPriceItineraryTimeout)

	headers := map[string]string{
		"Action":    string(GetPriceItinerary),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp GetPriceItineraryResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// SellJourneys sells the selected journeys
func (c *Client) SellJourneys(ctx context.Context, req SellJourneysRequest) (*SellJourneysResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "SellJourneys")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SellJourneys request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(sellJourneysTimeout)

	headers := map[string]string{
		"Action":    string(SellJourneys),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp SellJourneysResponse
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
