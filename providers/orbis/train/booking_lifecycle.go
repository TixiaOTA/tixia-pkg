package train

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

const issueBookingTimeout = 30 * time.Second
const cancelBookingTimeout = 30 * time.Second

// IssueBooking issues a previously committed booking. It never retries an
// ambiguous provider outcome; callers can reconcile with GetBooking.
func (c *clientImpl) IssueBooking(ctx context.Context, req IssueBookingRequest) (*IssueBookingResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "IssueBooking")
	defer span.End()
	if err := req.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	signin, err := c.signIn(ctx)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{"Action": "issued", "Signature": signin.Result.Response.Signature}
	timeout := issueBookingTimeout
	noRetry := client.RetryConfig{MaxRetries: 0}
	httpResp, sendErr := c.sendRequest(ctx, req, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout, Retry: &noRetry})
	if sendErr != nil {
		return nil, newProviderError(ErrorInternal, sendErr.Error())
	}
	var resp IssueBookingResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	if resp.Result.Response.OrderNumber != req.OrderNumber {
		return &resp, newProviderError(ErrorInternal, "Orbis issued returned different order number")
	}
	if !strings.EqualFold(resp.Result.Response.State, "issued") {
		return &resp, newProviderError(ErrorInternal, fmt.Sprintf("Orbis issued returned state %q", resp.Result.Response.State))
	}
	for _, booking := range resp.Result.Response.ProviderBookings {
		if !strings.EqualFold(booking.State, "issued") {
			return &resp, newProviderError(ErrorInternal, fmt.Sprintf("Orbis issued returned provider booking state %q", booking.State))
		}
	}
	if len(resp.Result.Response.ProviderBookings) == 0 {
		return &resp, newProviderError(ErrorInternal, "Orbis issued returned no provider bookings")
	}
	return &resp, nil
}

// CancelBooking cancels a held booking before issuance. A partial provider
// result is returned alongside the error for caller-side reconciliation.
func (c *clientImpl) CancelBooking(ctx context.Context, req CancelBookingRequest) (*CancelBookingResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "CancelBooking")
	defer span.End()
	if err := req.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	signin, err := c.signIn(ctx)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{"Action": "cancel_booking", "Signature": signin.Result.Response.Signature}
	timeout := cancelBookingTimeout
	noRetry := client.RetryConfig{MaxRetries: 0}
	httpResp, sendErr := c.sendRequest(ctx, req, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout, Retry: &noRetry})
	if sendErr != nil {
		return nil, newProviderError(ErrorInternal, sendErr.Error())
	}
	var resp CancelBookingResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	if len(resp.Result.Response) == 0 {
		return &resp, newProviderError(ErrorInternal, "Orbis cancel booking returned no provider results")
	}
	for i, result := range resp.Result.Response {
		if result.ErrorCode != ErrorSuccess || result.ErrorMsg != "" {
			return &resp, newProviderError(result.ErrorCode, fmt.Sprintf("Orbis cancel booking provider %d: %s", i+1, result.ErrorMsg))
		}
		if !strings.EqualFold(strings.TrimSpace(result.Status), "SUCCESS") {
			return &resp, newProviderError(ErrorInternal, fmt.Sprintf("Orbis cancel booking provider %d returned status %q", i+1, result.Status))
		}
	}
	return &resp, nil
}
