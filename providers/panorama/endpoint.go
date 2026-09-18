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
	signInTimeout            = 30 * time.Second
	searchTimeout            = 120 * time.Second
	getBalanceTimeout        = 30 * time.Second
	getBookingTimeout        = 30 * time.Second
	cancelTimeout            = 30 * time.Second
	getPriceItineraryTimeout = 30 * time.Second
	sellJourneysTimeout      = 30 * time.Second
	updateContactsTimeout    = 30 * time.Second
	updatePassengersTimeout  = 30 * time.Second
	commitBookingTimeout     = 30 * time.Second
	issuedTimeout            = 30 * time.Second

	getFFAvailabilityTimeout   = 30 * time.Second
	getSSRAvailabilityTimeout  = 30 * time.Second
	getSeatAvailabilityTimeout = 30 * time.Second
	sellSSRsTimeout            = 30 * time.Second
	assignSeatsTimeout         = 30 * time.Second
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

// Search searches for available flights.
func (c *Client) Search(ctx context.Context, req SearchRequest, options ...RequestOption) (*SearchResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "Search")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting Search request", zap.Any("request", sanitizedReq))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	payload := req
	payload.Signature = ""
	if payload.CarrierCodes == nil {
		payload.CarrierCodes = []string{}
	}
	if payload.PromoCodes == nil {
		payload.PromoCodes = []string{}
	}

	headers := map[string]string{
		"Action":    string(Search),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, payload, bookingAirlinePath, &headers, buildRequestOptions(searchTimeout, options...))
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp SearchResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if panoramaErr, ok := err.(*PanoramaError); ok {
			return nil, panoramaErr
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
func (c *Client) SellJourneys(ctx context.Context, req SellJourneysRequest, options ...RequestOption) (*SellJourneysResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "SellJourneys")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SellJourneys request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	headers := map[string]string{
		"Action":    string(SellJourneys),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers, buildRequestOptions(sellJourneysTimeout, options...))
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

// UpdateContacts updates booker and contact persons for a booking session
func (c *Client) UpdateContacts(ctx context.Context, req UpdateContactsRequest, options ...RequestOption) (*UpdateContactsResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "UpdateContacts")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting UpdateContacts request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	headers := map[string]string{
		"Action":    string(UpdateContacts),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers, buildRequestOptions(updateContactsTimeout, options...))
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp UpdateContactsResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// UpdatePassengers updates passengers for a booking session
func (c *Client) UpdatePassengers(ctx context.Context, req UpdatePassengersRequest, options ...RequestOption) (*UpdatePassengersResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "UpdatePassengers")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting UpdatePassengers request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	headers := map[string]string{
		"Action":    string(UpdatePassengers),
		"Signature": req.Signature,
	}

	httpResp, err := c.sendRequest(ctx, req, bookingAirlinePath, &headers, buildRequestOptions(updatePassengersTimeout, options...))
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	var resp UpdatePassengersResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// CommitBooking commits the current booking session
func (c *Client) CommitBooking(ctx context.Context, req CommitBookingRequest) (*CommitBookingResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "CommitBooking")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CommitBooking request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(commitBookingTimeout)

	headers := map[string]string{
		"Action":    string(CommitBooking),
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

	var resp CommitBookingResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// Issued issues a booking by order number
func (c *Client) Issued(ctx context.Context, req IssuedRequest) (*IssuedResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "Issued")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting Issued request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(issuedTimeout)

	headers := map[string]string{
		"Action":    string(Issued),
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

	var resp IssuedResponse
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

// GetFFAvailability gets the frequent flyer program availability for a booking session
func (c *Client) GetFFAvailability(ctx context.Context, req GetFFAvailabilityRequest) (*GetFFAvailabilityResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetFFAvailability")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetFFAvailability request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(getFFAvailabilityTimeout)

	headers := map[string]string{
		"Action":    string(GetFFAvailability),
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

	var resp GetFFAvailabilityResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// GetSSRAvailability gets the sellable SSR (ancillary service) availability for a booking session
func (c *Client) GetSSRAvailability(ctx context.Context, req GetSSRAvailabilityRequest) (*GetSSRAvailabilityResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetSSRAvailability")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetSSRAvailability request", zap.Any("request", sanitizedReq))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(getSSRAvailabilityTimeout)

	headers := map[string]string{
		"Action":    string(GetSSRAvailability),
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

	var resp GetSSRAvailabilityResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// GetSeatAvailability gets the seat map availability for a booking session
func (c *Client) GetSeatAvailability(ctx context.Context, req GetSeatAvailabilityRequest) (*GetSeatAvailabilityResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "GetSeatAvailability")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetSeatAvailability request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(getSeatAvailabilityTimeout)

	headers := map[string]string{
		"Action":    string(GetSeatAvailability),
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

	var resp GetSeatAvailabilityResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// SellSSRs sells the selected SSRs (ancillary services) for a booking session
func (c *Client) SellSSRs(ctx context.Context, req SellSSRsRequest) (*SellSSRsResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "SellSSRs")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SellSSRs request", zap.Any("request", sanitizedReq))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(sellSSRsTimeout)

	headers := map[string]string{
		"Action":    string(SellSSRs),
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

	var resp SellSSRsResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}

// AssignSeats assigns the selected seats for a booking session
func (c *Client) AssignSeats(ctx context.Context, req AssignSeatsRequest) (*AssignSeatsResponse, *PanoramaError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "AssignSeats")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting AssignSeats request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	c.httpClient.SetTimeout(assignSeatsTimeout)

	headers := map[string]string{
		"Action":    string(AssignSeats),
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

	var resp AssignSeatsResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if _, ok := err.(*PanoramaError); ok {
			return nil, err.(*PanoramaError)
		}
		return nil, buildPanoramaError(InternalServerError, err.Error())
	}

	return &resp, nil
}
