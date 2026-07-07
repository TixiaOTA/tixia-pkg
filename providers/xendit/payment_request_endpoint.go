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
	createPaymentRequestPath    = "/v3/payment_requests"
	createPaymentRequestTimeout = 15 * time.Second
)

// CreatePaymentRequest initiates a direct channel payment using Xendit Payments API.
func (c *Client) CreatePaymentRequest(ctx context.Context, req CreatePaymentRequestRequest) (*CreatePaymentRequestResponse, *XenditError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "CreatePaymentRequest")
	defer span.End()

	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CreatePaymentRequest request", zap.Any("request", sanitizedReq))

	if err := req.Validate(); err != nil {
		return nil, &XenditError{
			ErrorCode: InvalidRequest,
			Message:   err.Error(),
		}
	}

	c.httpClient.SetTimeout(createPaymentRequestTimeout)

	httpResp, err := c.sendRequest(ctx, req, createPaymentRequestPath, nil)
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
