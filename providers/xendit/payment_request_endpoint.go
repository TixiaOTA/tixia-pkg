package xendit

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

const (
	createPaymentRequestPath      = "/v3/payment_requests"
	paymentRequestAPIVersion      = "2024-11-11"
	createPaymentRequestTimeout   = 15 * time.Second
	getPaymentRequestTimeout      = 15 * time.Second
	cancelPaymentRequestTimeout   = 15 * time.Second
	simulatePaymentRequestTimeout = 15 * time.Second
)

var paymentRequestHeaders = map[string]string{
	"api-version": paymentRequestAPIVersion,
}

// CreatePaymentRequest initiates a direct channel payment using Xendit Payments API.
func (c *Client) CreatePaymentRequest(ctx context.Context, req CreatePaymentRequestRequest) (*CreatePaymentRequestResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CreatePaymentRequest")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CreatePaymentRequest request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(createPaymentRequestTimeout)

	httpResp, err := c.sendRequest(ctx, req, createPaymentRequestPath, &paymentRequestHeaders)
	if err != nil {
		return nil, toXenditError(err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	var resp CreatePaymentRequestResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CreatePaymentRequest completed successfully",
		zap.String("payment_request_id", resp.PaymentRequestID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// GetPaymentRequest retrieves the status and details of a payment request.
func (c *Client) GetPaymentRequest(ctx context.Context, req GetPaymentRequestStatusRequest) (*CreatePaymentRequestResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "GetPaymentRequest")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetPaymentRequest request",
		zap.String("payment_request_id", req.PaymentRequestID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(getPaymentRequestTimeout)

	path := fmt.Sprintf("%s/%s", createPaymentRequestPath, strings.TrimSpace(req.PaymentRequestID))
	httpResp, err := c.sendGetRequest(ctx, path, &paymentRequestHeaders)
	if err != nil {
		return nil, toXenditError(err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	var resp CreatePaymentRequestResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("GetPaymentRequest completed successfully",
		zap.String("payment_request_id", resp.PaymentRequestID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// CancelPaymentRequest cancels a payment request and prevents the end user from completing payment.
func (c *Client) CancelPaymentRequest(ctx context.Context, req CancelPaymentRequestStatusRequest) (*CreatePaymentRequestResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CancelPaymentRequest")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CancelPaymentRequest request",
		zap.String("payment_request_id", req.PaymentRequestID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(cancelPaymentRequestTimeout)

	path := fmt.Sprintf("%s/%s/cancel", createPaymentRequestPath, strings.TrimSpace(req.PaymentRequestID))
	httpResp, err := c.sendPostRequest(ctx, path, nil, &paymentRequestHeaders)
	if err != nil {
		return nil, toXenditError(err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	var resp CreatePaymentRequestResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CancelPaymentRequest completed successfully",
		zap.String("payment_request_id", resp.PaymentRequestID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// SimulatePaymentRequest simulates payment completion for a payment request in test mode.
func (c *Client) SimulatePaymentRequest(ctx context.Context, req SimulatePaymentRequestRequest) (*SimulatePaymentRequestResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "SimulatePaymentRequest")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SimulatePaymentRequest request", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(simulatePaymentRequestTimeout)

	path := fmt.Sprintf("%s/%s/simulate", createPaymentRequestPath, strings.TrimSpace(req.PaymentRequestID))
	httpResp, err := c.sendPostRequest(ctx, path, req, &paymentRequestHeaders)
	if err != nil {
		return nil, toXenditError(err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	var resp SimulatePaymentRequestResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("SimulatePaymentRequest completed successfully",
		zap.String("payment_request_id", req.PaymentRequestID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}
