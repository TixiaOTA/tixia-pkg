package train

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/providers/panorama"
)

const bookingTrainPath = "/booking/train"

type Config struct {
	BaseURL     string
	RetryConfig *client.RetryConfig
}

type Client struct {
	httpClient *client.HTTPClient
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	httpConfig := client.HTTPConfig{
		BaseURL: strings.TrimRight(config.BaseURL, "/"),
		Retry:   config.RetryConfig,
	}
	return &Client{httpClient: client.NewHTTPClient(httpConfig)}, nil
}

// Search uses the signature from panorama.Client.SignIn.
// Call SignIn before each Train operation, just as in the Flight flow.
func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, *panorama.PanoramaError) {
	if err := request.Validate(); err != nil {
		return nil, providerError(panorama.InvalidRequest, err.Error())
	}
	if c == nil || c.httpClient == nil {
		return nil, providerError(panorama.InternalServerError, "train client is not initialized")
	}

	headers := map[string]string{
		"Action":    "search",
		"Signature": request.Signature,
		"Accept":    "application/json",
	}
	timeout := 120 * time.Second
	resp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    bookingTrainPath,
		Body:    request,
		Headers: &headers,
		Options: &client.RequestOptions{Timeout: &timeout},
	})
	if err != nil {
		return nil, providerError(panorama.InternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, providerError(panorama.InternalServerError, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, providerError(panorama.InternalServerError, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
	}

	// Orbis error payloads may contain response as a string, rather than an object.
	var envelope struct {
		Result struct {
			ErrorCode panorama.PanoramaErrorCode `json:"error_code"`
			ErrorMsg  string                    `json:"error_msg"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, providerError(panorama.InternalServerError, err.Error())
	}
	if envelope.Result.ErrorCode != panorama.Success || envelope.Result.ErrorMsg != "" {
		return nil, providerError(envelope.Result.ErrorCode, envelope.Result.ErrorMsg)
	}
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, providerError(panorama.InternalServerError, err.Error())
	}
	return &result, nil
}

func providerError(code panorama.PanoramaErrorCode, message string) *panorama.PanoramaError {
	return &panorama.PanoramaError{
		Result: panorama.PanoramaErrorResult{
			BasicPanoramaResult: panorama.BasicPanoramaResult{
				ErrorCode: code,
				ErrorMsg:  message,
			},
		},
	}
}
