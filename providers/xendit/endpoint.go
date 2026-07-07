package xendit

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

const (
	createSessionPath     = "/sessions"
	createSessionTimeout  = 15 * time.Second
	getSessionTimeout     = 15 * time.Second
	cancelSessionTimeout  = 15 * time.Second
	getPaymentTimeout     = 15 * time.Second
	cancelPaymentTimeout  = 15 * time.Second
	capturePaymentTimeout = 15 * time.Second
)

// CreateSession initiates a one-time payment session using Xendit hosted checkout.
func (c *Client) CreateSession(ctx context.Context, req CreateSessionRequest) (*CreateSessionResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CreateSession")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CreateSession request", zap.Any("request", sanitizedReq))

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(createSessionTimeout)

	httpResp, err := c.sendRequest(ctx, req, createSessionPath, nil)
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

	var resp CreateSessionResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CreateSession completed successfully",
		zap.String("payment_session_id", resp.PaymentSessionID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// GetSession retrieves the status of a payment session and related payment references.
func (c *Client) GetSession(ctx context.Context, req GetSessionRequest) (*GetSessionResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "GetSession")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetSession request",
		zap.String("payment_session_id", req.PaymentSessionID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(getSessionTimeout)

	path := fmt.Sprintf("/sessions/%s", strings.TrimSpace(req.PaymentSessionID))
	httpResp, err := c.sendGetRequest(ctx, path, nil)
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

	var resp GetSessionResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("GetSession completed successfully",
		zap.String("payment_session_id", resp.PaymentSessionID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// CancelSession cancels an active payment session and invalidates its hosted checkout URL.
func (c *Client) CancelSession(ctx context.Context, req CancelSessionRequest) (*CancelSessionResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CancelSession")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CancelSession request",
		zap.String("payment_session_id", req.PaymentSessionID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(cancelSessionTimeout)

	path := fmt.Sprintf("/sessions/%s/cancel", strings.TrimSpace(req.PaymentSessionID))
	httpResp, err := c.sendPostRequest(ctx, path, nil, nil)
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

	var resp CancelSessionResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CancelSession completed successfully",
		zap.String("payment_session_id", resp.PaymentSessionID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// GetPayment retrieves the status and details of a payment.
func (c *Client) GetPayment(ctx context.Context, req GetPaymentRequest) (*GetPaymentResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "GetPayment")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetPayment request",
		zap.String("payment_id", req.PaymentID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(getPaymentTimeout)

	path := fmt.Sprintf("/v3/payments/%s", strings.TrimSpace(req.PaymentID))
	httpResp, err := c.sendGetRequest(ctx, path, nil)
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

	var resp GetPaymentResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("GetPayment completed successfully",
		zap.String("payment_id", resp.PaymentID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// CancelPayment cancels an authorized payment before capture.
func (c *Client) CancelPayment(ctx context.Context, req CancelPaymentRequest) (*CancelPaymentResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CancelPayment")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CancelPayment request",
		zap.String("payment_id", req.PaymentID),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(cancelPaymentTimeout)

	path := fmt.Sprintf("/v3/payments/%s/cancel", strings.TrimSpace(req.PaymentID))
	httpResp, err := c.sendPostRequest(ctx, path, nil, nil)
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

	var resp CancelPaymentResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CancelPayment completed successfully",
		zap.String("payment_id", resp.PaymentID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}

// CapturePayment captures an authorized payment to complete collection.
func (c *Client) CapturePayment(ctx context.Context, req CapturePaymentRequest) (*CapturePaymentResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CapturePayment")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CapturePayment request",
		zap.String("payment_id", req.PaymentID),
		zap.Any("request", sanitizedReq),
	)

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(capturePaymentTimeout)

	path := fmt.Sprintf("/v3/payments/%s/capture", strings.TrimSpace(req.PaymentID))
	httpResp, err := c.sendPostRequest(ctx, path, req, nil)
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

	var resp CapturePaymentResponse
	if err := parseResponseBody(bodyBytes, &resp); err != nil {
		if xerr, ok := err.(*XenditError); ok {
			return nil, xerr
		}

		return nil, &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CapturePayment completed successfully",
		zap.String("payment_id", resp.PaymentID),
		zap.String("status", resp.Status),
	)

	return &resp, nil
}
