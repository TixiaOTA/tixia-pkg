package panorama

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

func buildHeaders(additionalHeaders *map[string]string) map[string]string {
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

	headersMap := buildHeaders(headers)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headersMap,
	})
	if err != nil {
		sanitizedReq, _ := logger.Sanitize(req)
		c.logger.WithTrace(ctx, currentSpanId).Error("Error sending request: ",
			zap.Any("request", sanitizedReq),
			zap.Any("error", err),
		)

		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: err.Error(),
		}
	}

	// If the response is not 200 OK, return an API error
	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()

		// Try to unmarshal as PanoramaError first
		var panoramaError PanoramaError
		if err := json.Unmarshal(bodyBytes, &panoramaError); err == nil {
			// Check if it's a valid PanoramaError
			if panoramaError.Result.ErrorCode != Success || panoramaError.Result.ErrorMsg != "" {
				return nil, &PanoramaError{
					Result: PanoramaErrorResult{
						BasicPanoramaResult: BasicPanoramaResult{
							ErrorCode: panoramaError.Result.ErrorCode,
							ErrorMsg:  panoramaError.Result.ErrorMsg,
						},
						Response:           panoramaError.Result.Response,
						ErrorAdditionalMsg: panoramaError.Result.ErrorAdditionalMsg,
					},
				}
			}
		}

		return nil, &APIError{
			Code:    httpResp.StatusCode,
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

// parseResponseBody unmarshals the response body into the given response.
// If the response is a PanoramaError, it will be unmarshaled into the PanoramaError.
// If the response is not a PanoramaError, it will be unmarshaled into the response.
func parseResponseBody(bodyBytes []byte, response any) error {
	// Try to unmarshal as PanoramaError first
	var panoramaError PanoramaError
	if err := json.Unmarshal(bodyBytes, &panoramaError); err == nil {
		// Check if it's a valid PanoramaError
		if panoramaError.Result.ErrorCode != Success || panoramaError.Result.ErrorMsg != "" {
			return &PanoramaError{
				Result: PanoramaErrorResult{
					BasicPanoramaResult: BasicPanoramaResult{
						ErrorCode: panoramaError.Result.ErrorCode,
						ErrorMsg:  panoramaError.Result.ErrorMsg,
					},
					Response:           panoramaError.Result.Response,
					ErrorAdditionalMsg: panoramaError.Result.ErrorAdditionalMsg,
				},
			}
		}
	}

	// If not a PanoramaError, unmarshal as response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return err
	}

	return nil
}
