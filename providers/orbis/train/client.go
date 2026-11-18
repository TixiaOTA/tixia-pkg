package train

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

const (
	sessionPath      = "/session"
	bookingTrainPath = "/booking/train"
	signInTimeout    = 30 * time.Second
	searchTimeout    = 120 * time.Second
	seatAvailabilityTimeout = 30 * time.Second
	commitBookingTimeout = 30 * time.Second
)

type Config struct {
	BaseURL     string
	Credentials Credentials
	RetryConfig *client.RetryConfig
}

type Client interface {
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError)
	GetSeatAvailability(ctx context.Context, signature string, request GetSeatAvailabilityRequest) (*GetSeatAvailabilityResponse, *ProviderError)
	CommitBooking(ctx context.Context, signature string, request CommitBookingRequest) (*CommitBookingResponse, *ProviderError)
}

type clientImpl struct {
	httpClient  *client.HTTPClient
	credentials Credentials
}

var _ Client = (*clientImpl)(nil)

func New(config Config) (Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}
	if err := config.Credentials.Validate(); err != nil {
		return nil, err
	}
	if config.RetryConfig == nil {
		config.RetryConfig = &client.RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    5 * time.Second,
			MaxBackoff:        15 * time.Second,
			BackoffMultiplier: 2,
		}
	}
	return &clientImpl{
		httpClient: client.NewHTTPClient(client.HTTPConfig{
			BaseURL: strings.TrimRight(config.BaseURL, "/"),
			Retry:   config.RetryConfig,
		}),
		credentials: config.Credentials,
	}, nil
}

// signIn is internal; each search starts a fresh provider session.
func (c *clientImpl) signIn(ctx context.Context) (*signInResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "SignIn")
	defer span.End()

	headers := map[string]string{"Action": "signin"}
	timeout := signInTimeout
	httpResp, err := c.sendRequest(ctx, c.credentials, sessionPath, &headers, client.RequestOptions{Timeout: &timeout})
	if err != nil {
		return nil, newProviderError(ErrorInternal, err.Error())
	}

	var resp signInResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	if strings.TrimSpace(resp.Result.Response.Signature) == "" {
		return nil, newProviderError(ErrorInternal, "Orbis signin returned empty signature")
	}
	return &resp, nil
}

// Search signs in before calling the provider, and returns the signature for
// the backend to store for subsequent operations in the same booking session.
func (c *clientImpl) Search(ctx context.Context, req SearchRequest) (*SearchResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "Search")
	defer span.End()

	if err := req.Validate(); err != nil {
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}

	signin, err := c.signIn(ctx)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Action":    "search",
		"Signature": signin.Result.Response.Signature,
	}
	timeout := searchTimeout
	httpResp, sendErr := c.sendRequest(ctx, req, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout})
	if sendErr != nil {
		return nil, newProviderError(ErrorInternal, sendErr.Error())
	}

	var resp SearchResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	resp.Signature = signin.Result.Response.Signature
	return &resp, nil
}
