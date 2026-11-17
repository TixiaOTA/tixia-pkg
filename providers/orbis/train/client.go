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
)

const (
	sessionPath      = "/session"
	bookingTrainPath = "/booking/train"
)

type Config struct {
	BaseURL     string
	Credentials Credentials
	RetryConfig *client.RetryConfig
}

type Client struct {
	httpClient *client.HTTPClient
	credentials Credentials
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if err := config.Credentials.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		httpClient: client.NewHTTPClient(client.HTTPConfig{
			BaseURL: strings.TrimRight(config.BaseURL, "/"),
			Retry: config.RetryConfig,
		}),
		credentials: config.Credentials,
	}, nil
}

// SignIn creates a provider session. Call once before each independent Train flow.
func (c *Client) SignIn(ctx context.Context) (*SignInResponse, *ProviderError) {
	if c == nil || c.httpClient == nil {
		return nil, newProviderError(ErrorInternal, "train client is not initialized")
	}
	var response SignInResponse
	if err := c.call(ctx, sessionPath, "signin", "", c.credentials, 30*time.Second, &response); err != nil {
		return nil, err
	}
	if strings.TrimSpace(response.Result.Response.Signature) == "" {
		return nil, newProviderError(ErrorInternal, "Orbis signin returned empty signature")
	}
	return &response, nil
}

// Search uses a session signature obtained from SignIn. The signature is only
// sent in the header, never in the JSON request body.
func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError) {
	if err := request.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	if c == nil || c.httpClient == nil {
		return nil, newProviderError(ErrorInternal, "train client is not initialized")
	}
	var response SearchResponse
	if err := c.call(ctx, bookingTrainPath, "search", request.Signature, request, 120*time.Second, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Callers start each flow with SignIn, then reuse its signature for the
// subsequent provider calls belonging to that flow.
func (c *Client) call(ctx context.Context, path, action, signature string, body any, timeout time.Duration, dest any) *ProviderError {
	headers := map[string]string{"Action": action, "Accept": "application/json"}
	if signature != "" {
		headers["Signature"] = signature
	}
	resp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path: path, Body: body, Headers: &headers,
		Options: &client.RequestOptions{Timeout: &timeout},
	})
	if err != nil {
		return newProviderError(ErrorInternal, err.Error())
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return newProviderError(ErrorInternal, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return newProviderError(ErrorInternal, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(payload)))
	}
	var envelope struct {
		Result BasicResult `json:"result"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return newProviderError(ErrorInternal, err.Error())
	}
	if envelope.Result.ErrorCode != ErrorSuccess || envelope.Result.ErrorMsg != "" {
		return newProviderError(envelope.Result.ErrorCode, envelope.Result.ErrorMsg)
	}
	if err := json.Unmarshal(payload, dest); err != nil {
		return newProviderError(ErrorInternal, err.Error())
	}
	return nil
}
