package train

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/attribute"
)

// sendRequest follows the Panorama HTTP helper pattern. The caller handles
// lifecycle logging so successful provider calls do not produce duplicate logs.
func (c *clientImpl) sendRequest(ctx context.Context, req any, path string, extra *map[string]string, options ...client.RequestOptions) (*http.Response, error) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "sendRequest")
	defer span.End()

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	if extra != nil {
		maps.Copy(headers, *extra)
	}
	span.SetAttributes(
		attribute.String("provider.action", headers["Action"]),
		attribute.String("http.route", path),
	)

	post := client.PostJSONRequest{Path: path, Body: req, Headers: &headers}
	if len(options) > 0 {
		post.Options = &options[0]
	}
	resp, err := c.httpClient.PostJSON(ctx, post)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "provider HTTP request failed")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		// Do not include the raw body: provider errors may contain sensitive data.
		err := fmt.Errorf("orbis train HTTP status %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "provider returned non-200 status")
		return nil, err
	}
	return resp, nil
}

// decodeResponse reads and decodes the provider response for both SignIn and
// Search. It owns closing the response body.
func decodeResponse(resp *http.Response, destination any) *ProviderError {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newProviderError(ErrorInternal, fmt.Sprintf("read Orbis response: %v", err))
	}
	return parseResponseBody(body, destination)
}

func parseResponseBody(body []byte, destination any) *ProviderError {
	// Provider errors may use response: "" rather than an object.
	var status struct {
		Result BasicResult `json:"result"`
	}
	if err := json.Unmarshal(body, &status); err != nil {
		return newProviderError(ErrorInternal, fmt.Sprintf("invalid Orbis response: %v", err))
	}
	if status.Result.ErrorCode != ErrorSuccess || status.Result.ErrorMsg != "" {
		return newProviderError(status.Result.ErrorCode, status.Result.ErrorMsg)
	}
	if err := json.Unmarshal(body, destination); err != nil {
		return newProviderError(ErrorInternal, fmt.Sprintf("invalid Orbis payload: %v", err))
	}
	return nil
}
