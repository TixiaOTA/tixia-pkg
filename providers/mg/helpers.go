package mg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

func buildHeaders(additionalHeaders *map[string]string) map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	if additionalHeaders != nil {
		for key, value := range *additionalHeaders {
			headers[key] = value
		}
	}

	return headers
}

func (c *Client) sendRequest(ctx context.Context, req any, path string, headers *map[string]string) (*http.Response, error) {
	headersMap := buildHeaders(headers)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headersMap,
	})
	if err != nil {
		return nil, err
	}

	// If the response is not 200 OK, return an API error
	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	return httpResp, nil
}

// parseResponseBody unmarshals the response body into the given response.
// If the response is a MGError, it will be unmarshaled into the MGError.
// If the response is not a MGError, it will be unmarshaled into the response.
func parseResponseBody(bodyBytes []byte, response any) error {
	// Try to unmarshal as MGError first
	var mgError MGError
	if err := json.Unmarshal(bodyBytes, &mgError); err == nil {
		// Check if it's a valid MGError
		if !mgError.Status || mgError.ErrorCode != "" || mgError.ErrorMessage != "" {
			return &MGError{
				Status:       false,
				ErrorCode:    mgError.ErrorCode,
				ErrorMessage: mgError.ErrorMessage,
				RawBody:      &bodyBytes,
			}
		}
	}

	// If not a MGError, unmarshal as response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return err
	}

	return nil
}
