package mg

import (
	"context"
	"io"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

func (c *Client) buildHeaders(additionalHeaders *map[string]string) map[string]string {
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
	headersMap := c.buildHeaders(headers)

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
