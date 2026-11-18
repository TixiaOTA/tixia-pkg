package train

import (
	"context"
	"strings"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

// GetSeatAvailability fetches provider cabin layouts for one or two fare codes.
// Each outer response element corresponds to the same-index request fare code.
func (c *clientImpl) GetSeatAvailability(ctx context.Context, signature string, req GetSeatAvailabilityRequest) (*GetSeatAvailabilityResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "GetSeatAvailability")
	defer span.End()

	if strings.TrimSpace(signature) == "" {
		return nil, newProviderError(ErrorInvalidInput, "signature is required")
	}
	if err := req.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	headers := map[string]string{"Action": "get_seat_availability", "Signature": signature}
	timeout := seatAvailabilityTimeout
	httpResp, err := c.sendRequest(ctx, req, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout})
	if err != nil {
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	var resp GetSeatAvailabilityResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	if resp.Result.Response == nil || len(resp.Result.Response) != len(req.FareCodes) {
		return nil, newProviderError(ErrorInternal, "Orbis seat availability response does not match requested fares")
	}
	return &resp, nil
}

// CommitBooking holds a provider booking, without issuing tickets.
// Explicitly disable shared HTTP retries to avoid double commits after timeouts.
func (c *clientImpl) CommitBooking(ctx context.Context, signature string, req CommitBookingRequest) (*CommitBookingResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "CommitBooking")
	defer span.End()

	if strings.TrimSpace(signature) == "" {
		return nil, newProviderError(ErrorInvalidInput, "signature is required")
	}
	if err := req.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	headers := map[string]string{"Action": "commit_booking", "Signature": signature}
	timeout := commitBookingTimeout
	noRetry := client.RetryConfig{MaxRetries: 0}
	httpResp, err := c.sendRequest(ctx, req, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout, Retry: &noRetry})
	if err != nil {
		// The backend must reconcile an ambiguous timeout; never retry automatically.
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	var resp CommitBookingResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	if strings.TrimSpace(resp.Result.Response.OrderNumber) == "" {
		return nil, newProviderError(ErrorInternal, "Orbis commit booking returned empty order number")
	}
	return &resp, nil
}
