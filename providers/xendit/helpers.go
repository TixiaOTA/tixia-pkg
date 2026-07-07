package xendit

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

func (c *Client) buildHeaders(additionalHeaders *map[string]string) map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": c.authHeader,
	}

	if c.config.ForUserID != "" {
		headers["for-user-id"] = c.config.ForUserID
	}

	if c.config.SplitRuleID != "" {
		headers["with-split-rule"] = c.config.SplitRuleID
	}

	if additionalHeaders != nil {
		maps.Copy(headers, *additionalHeaders)
	}

	return headers
}

func (c *Client) sendRequest(ctx context.Context, req any, path string, headers *map[string]string) (*http.Response, error) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "sendRequest")
	defer span.End()

	headersMap := c.buildHeaders(headers)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headersMap,
	})
	if err != nil {
		sanitizedReq, _ := logger.Sanitize(req)
		c.logger.WithTrace(ctx, currentSpanId).Error("Error sending request",
			zap.Any("request", sanitizedReq),
			zap.Error(err),
		)

		return nil, &APIError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if httpResp.StatusCode != http.StatusOK && httpResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		c.logger.WithTrace(ctx, currentSpanId).Error("HTTP request error",
			zap.Int("status_code", httpResp.StatusCode),
			zap.ByteString("body", bodyBytes),
		)

		var xenditErr XenditError
		if err := json.Unmarshal(bodyBytes, &xenditErr); err != nil || xenditErr.ErrorCode == "" {
			return nil, &XenditError{
				ErrorCode:  InternalServerError,
				Message:    string(bodyBytes),
				HTTPStatus: httpResp.StatusCode,
			}
		}

		xenditErr.HTTPStatus = httpResp.StatusCode
		return nil, &xenditErr
	}

	return httpResp, nil
}

func (c *Client) sendGetRequest(ctx context.Context, path string, headers *map[string]string) (*http.Response, error) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "sendGetRequest")
	defer span.End()

	headersMap := c.buildHeaders(headers)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headersMap,
	})
	if err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("Error sending GET request",
			zap.String("path", path),
			zap.Error(err),
		)

		return nil, &APIError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		c.logger.WithTrace(ctx, currentSpanId).Error("HTTP GET request error",
			zap.String("path", path),
			zap.Int("status_code", httpResp.StatusCode),
			zap.ByteString("body", bodyBytes),
		)

		var xenditErr XenditError
		if err := json.Unmarshal(bodyBytes, &xenditErr); err != nil || xenditErr.ErrorCode == "" {
			return nil, &XenditError{
				ErrorCode:  InternalServerError,
				Message:    string(bodyBytes),
				HTTPStatus: httpResp.StatusCode,
			}
		}

		xenditErr.HTTPStatus = httpResp.StatusCode
		return nil, &xenditErr
	}

	return httpResp, nil
}

func (c *Client) sendPostRequest(ctx context.Context, path string, body any, headers *map[string]string) (*http.Response, error) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "xendit", "sendPostRequest")
	defer span.End()

	if body == nil {
		body = struct{}{}
	}

	headersMap := c.buildHeaders(headers)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    body,
		Headers: &headersMap,
	})
	if err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("Error sending POST request",
			zap.String("path", path),
			zap.Error(err),
		)

		return nil, &APIError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		c.logger.WithTrace(ctx, currentSpanId).Error("HTTP POST request error",
			zap.String("path", path),
			zap.Int("status_code", httpResp.StatusCode),
			zap.ByteString("body", bodyBytes),
		)

		var xenditErr XenditError
		if err := json.Unmarshal(bodyBytes, &xenditErr); err != nil || xenditErr.ErrorCode == "" {
			return nil, &XenditError{
				ErrorCode:  InternalServerError,
				Message:    string(bodyBytes),
				HTTPStatus: httpResp.StatusCode,
			}
		}

		xenditErr.HTTPStatus = httpResp.StatusCode
		return nil, &xenditErr
	}

	return httpResp, nil
}

func parseResponseBody(bodyBytes []byte, response any) error {
	if err := json.Unmarshal(bodyBytes, response); err != nil {
		return &XenditError{
			ErrorCode: InternalServerError,
			Message:   err.Error(),
		}
	}

	return nil
}

func toXenditError(err error) *XenditError {
	if err == nil {
		return nil
	}

	if xerr, ok := err.(*XenditError); ok {
		return xerr
	}

	if apiErr, ok := err.(*APIError); ok {
		return &XenditError{
			ErrorCode:  InternalServerError,
			Message:    apiErr.Message,
			HTTPStatus: apiErr.Code,
		}
	}

	return &XenditError{
		ErrorCode: InternalServerError,
		Message:   err.Error(),
	}
}
