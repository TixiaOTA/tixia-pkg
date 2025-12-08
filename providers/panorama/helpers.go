package panorama

import (
	"context"
	"io"
	"maps"
	"net/http"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

func (c *Client) buildHeaders(additionalHeaders *map[string]string) map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	if additionalHeaders != nil {
		maps.Copy(headers, *additionalHeaders)
	}

	return headers
}

func (c *Client) sendRequest(ctx context.Context, req any, path string, headers *map[string]string) (*http.Response, error) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "panorama", "sendRequest")
	defer span.End()

	headersMap := c.buildHeaders(headers)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headersMap,
	})
	if err != nil {
		c.logger.WithTrace(ctx, currentSpanId).Error("Error sending request: ",
			zap.Any("request", req),
			zap.Any("error", err),
		)

		return nil, &APIError{
			Code:    InternalServerError,
			Message: err.Error(),
		}
	}

	// If the response is not 200 OK, return an API error
	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		return nil, &APIError{
			Code:    InternalServerError,
			Message: string(bodyBytes),
		}
	}

	return httpResp, nil
}

// buildPanoramaError builds a PanoramaError with the given code and message.
// Use this for errors that are API errors or not Panorama errors.
func buildPanoramaError(code PanoramaErrorCode, message string) *PanoramaError {
	return &PanoramaError{
		Result: PanoramaErrorResult{
			BasicPanoramaResult: BasicPanoramaResult{
				ErrorCode: code,
				ErrorMsg:  message,
			},
			Response:           "",
			ErrorAdditionalMsg: "",
		},
	}
}
