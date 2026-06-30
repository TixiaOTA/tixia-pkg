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
	createSessionPath    = "/sessions"
	createSessionTimeout = 15 * time.Second
	getSessionTimeout    = 15 * time.Second
	cancelSessionTimeout = 15 * time.Second
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
	httpResp, err := c.sendPostRequest(ctx, path, nil)
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
