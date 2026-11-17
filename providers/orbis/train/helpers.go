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
)

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
	post := client.PostJSONRequest{Path: path, Body: req, Headers: &headers}
	if len(options) > 0 {
		post.Options = &options[0]
	}
	resp, err := c.httpClient.PostJSON(ctx, post)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("orbis train HTTP status %d", resp.StatusCode)
	}
	return resp, nil
}

// decodeResponse reads and closes the body, then parses provider errors or data.
func decodeResponse(resp *http.Response, destination any) *ProviderError {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newProviderError(ErrorInternal, fmt.Sprintf("read Orbis response: %v", err))
	}
	return parseResponseBody(body, destination)
}

func parseResponseBody(body []byte, destination any) *ProviderError {
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
