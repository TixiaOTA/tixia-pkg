package train

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

// sendRequest mirrors the existing Panorama HTTP helper without coupling the
// Train package to Panorama types.
func (c *Client) sendRequest(ctx context.Context, req any, path string, extra *map[string]string, options ...client.RequestOptions) (*http.Response, error) {
	headers := map[string]string{"Content-Type": "application/json", "Accept": "application/json"}
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
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("orbis train HTTP %d: %s", resp.StatusCode, string(body))
	}
	return resp, nil
}

func parseResponseBody(body []byte, destination any) *ProviderError {
	// Provider errors may use response: "" instead of an object.
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
