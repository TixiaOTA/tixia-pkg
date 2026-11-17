package train

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
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

type Client interface {
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError)
}

type clientImpl struct {
	httpClient  *client.HTTPClient
	credentials Credentials
	logger      *logger.Logger
}

var _ Client = (*clientImpl)(nil)

func New(config Config, log *logger.Logger) (Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}
	if err := config.Credentials.Validate(); err != nil {
		return nil, err
	}
	if log == nil {
		return nil, errors.New("logger is required")
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
			Retry: config.RetryConfig,
		}),
		credentials: config.Credentials,
		logger: log,
	}, nil
}

// signIn is internal; each new search starts a fresh session.
func (c *clientImpl) signIn(ctx context.Context) (*signInResponse, *ProviderError) {
	ctx, span, _ := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "SignIn")
	defer span.End()
	headers := map[string]string{"Action": "signin"}
	timeout := signInTimeout
	resp, err := c.sendRequest(ctx, c.credentials, sessionPath, &headers, client.RequestOptions{Timeout: &timeout})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "signin failed")
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	var parsed signInResponse
	if decodeErr := decodeResponse(resp, &parsed); decodeErr != nil {
		span.RecordError(decodeErr)
		span.SetStatus(codes.Error, "signin response invalid")
		return nil, decodeErr
	}
	if strings.TrimSpace(parsed.Result.Response.Signature) == "" {
		err := newProviderError(ErrorInternal, "Orbis signin returned empty signature")
		span.RecordError(err)
		span.SetStatus(codes.Error, "signin missing signature")
		return nil, err
	}
	return &parsed, nil
}

// Search calls SignIn and Search in a single trace, logging one outcome.
func (c *clientImpl) Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "Search")
	defer span.End()
	start := time.Now()
	log := c.logger.WithTrace(ctx, spanID).With(zap.String("provider", "orbis"), zap.String("product", "train"))
	log.Info("Train search started", zap.String("direction", request.Direction), zap.Int("journey_count", len(request.JourneyList)))
	fail := func(step string, err *ProviderError) (*SearchResponse, *ProviderError) {
		span.RecordError(err)
		span.SetStatus(codes.Error, step+" failed")
		log.Error("Train search failed", zap.String("step", step), zap.Int("error_code", int(err.Code)), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
		return nil, err
	}

	if err := request.Validate(); err != nil {
		return fail("validate", newProviderError(ErrorInvalidInput, err.Error()))
	}
	signin, err := c.signIn(ctx)
	if err != nil {
		return fail("signin", err)
	}
	headers := map[string]string{"Action": "search", "Signature": signin.Result.Response.Signature}
	timeout := searchTimeout
	resp, sendErr := c.sendRequest(ctx, request, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout})
	if sendErr != nil {
		return fail("search", newProviderError(ErrorInternal, sendErr.Error()))
	}
	var result SearchResponse
	if decodeErr := decodeResponse(resp, &result); decodeErr != nil {
		return fail("decode", decodeErr)
	}
	result.Signature = signin.Result.Response.Signature
	log.Info("Train search completed", zap.Int("schedule_count", len(result.Result.Response.Schedules)), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
	return &result, nil
}
