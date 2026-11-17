package train

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

const (
	sessionPath      = "/session"
	bookingTrainPath = "/booking/train"
	signInTimeout    = 30 * time.Second
	searchTimeout    = 120 * time.Second
)

type Config struct {
	BaseURL     string
	Credentials Credentials
	RetryConfig *client.RetryConfig
}

// Client exposes Train provider operations without exposing session management.
type Client interface {
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError)
}

type clientImpl struct {
	httpClient  *client.HTTPClient
	credentials Credentials
}

var _ Client = (*clientImpl)(nil)

func New(config Config) (Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if err := config.Credentials.Validate(); err != nil {
		return nil, err
	}
	return &clientImpl{
		httpClient: client.NewHTTPClient(client.HTTPConfig{
			BaseURL: strings.TrimRight(config.BaseURL, "/"),
			Retry:   config.RetryConfig,
		}),
		credentials: config.Credentials,
	}, nil
}

// signIn stays internal to the Train package.
func (c *clientImpl) signIn(ctx context.Context) (*signInResponse, *ProviderError) {
	headers := map[string]string{"Action": "signin"}
	timeout := signInTimeout
	resp, err := c.sendRequest(ctx, c.credentials, sessionPath, &headers, client.RequestOptions{Timeout: &timeout})
	if err != nil {
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	var parsed signInResponse
	if err := parseResponseBody(body, &parsed); err != nil {
		return nil, err
	}
	if strings.TrimSpace(parsed.Result.Response.Signature) == "" {
		return nil, newProviderError(ErrorInternal, "Orbis signin returned empty signature")
	}
	return &parsed, nil
}

// Search signs in before every new search and keeps its signature in the
// internal response field for the backend to persist for the booking flow.
func (c *clientImpl) Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError) {
	if err := request.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	signin, err := c.signIn(ctx)
	if err != nil {
		return nil, err
	}
	signature := signin.Result.Response.Signature
	headers := map[string]string{
		"Action":    "search",
		"Signature": signature,
	}
	timeout := searchTimeout
	resp, sendErr := c.sendRequest(ctx, request, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout})
	if sendErr != nil {
		return nil, newProviderError(ErrorInternal, sendErr.Error())
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, newProviderError(ErrorInternal, readErr.Error())
	}
	var result SearchResponse
	if decodeErr := parseResponseBody(body, &result); decodeErr != nil {
		return nil, decodeErr
	}
	result.Signature = signature
	return &result, nil
}
