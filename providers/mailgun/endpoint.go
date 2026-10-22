package mailgun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"go.uber.org/zap"

	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

const messagesPathFormat = "/v3/%s/messages"

// SendEmail sends an email through the Mailgun API.
func (c *Client) SendEmail(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "provider", "mailgun", "SendEmail")
	defer span.End()

	if err := req.Validate(); err != nil {
		c.logger.WithTrace(ctx, spanID).Error("invalid send email request", zap.Error(err))
		return nil, err
	}

	path := fmt.Sprintf(messagesPathFormat, url.PathEscape(c.config.Domain))
	headers := map[string]string{
		"Authorization": basicAuthHeader(c.config.APIKey),
		"Accept":        "application/json",
	}

	c.logger.WithTrace(ctx, spanID).Info(
		"sending email",
		zap.String("from", req.From),
		zap.Int("recipients", len(req.To)+len(req.CC)+len(req.BCC)),
		zap.String("subject", req.Subject),
		zap.Int("attachments", len(req.Attachments)+len(req.Inline)),
	)

	httpResp, err := c.send(ctx, req, path, headers)
	if err != nil {
		c.logger.WithTrace(ctx, spanID).Error("failed to send email", zap.Error(err))
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("mailgun: read response: %w", err)
	}

	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices {
		apiError := parseAPIError(httpResp.StatusCode, body)
		c.logger.WithTrace(ctx, spanID).Error(
			"mailgun API error",
			zap.Int("status_code", httpResp.StatusCode),
			zap.String("message", apiError.Message),
			zap.ByteString("body", body),
		)
		return nil, apiError
	}

	var resp SendEmailResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("mailgun: decode response: %w", err)
	}

	return &resp, nil
}

func (c *Client) send(ctx context.Context, req SendEmailRequest, path string, headers map[string]string) (*http.Response, error) {
	if len(req.Attachments) > 0 || len(req.Inline) > 0 {
		body, contentType, err := buildMultipartBody(req)
		if err != nil {
			return nil, fmt.Errorf("mailgun: build multipart body: %w", err)
		}

		return c.httpClient.PostRaw(ctx, client.PostRawRequest{
			Path:        path,
			Body:        body,
			ContentType: contentType,
			Headers:     &headers,
		})
	}

	return c.httpClient.PostForm(ctx, client.PostFormRequest{
		Path:    path,
		Body:    buildForm(req),
		Headers: &headers,
	})
}
