package xendit

import (
	"context"
	"io"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

const (
	createSessionPath    = "/sessions"
	createSessionTimeout = 15 * time.Second
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
