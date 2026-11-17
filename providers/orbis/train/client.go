package train

import (
	"context"
	"errors"
	"fmt"
	"io"
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

// Client exposes provider operations; SignIn and session management stay private.
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
			Retry:   config.RetryConfig,
		}),
		credentials: config.Credentials,
		logger: log,
	}, nil
}

// signIn creates a provider session. It is called internally before Search.
func (c *clientImpl) signIn(ctx context.Context) (*signInResponse, *ProviderError) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "SignIn")
	defer span.End()
	start := time.Now()
	log := c.logger.WithTrace(ctx, spanID).With(
		zap.String("provider", "orbis"),
		zap.String("product", "train"),
		zap.String("operation", "signin"),
	)
	log.Info("Orbis train signin started", zap.String("stage", "signin.start"))

	headers := map[string]string{"Action": "signin"}
	timeout := signInTimeout
	resp, err := c.sendRequest(ctx, c.credentials, sessionPath, &headers, client.RequestOptions{Timeout: &timeout})
	if err != nil {
		span.SetStatus(codes.Error, "signin request failed")
		log.Error("Orbis train signin failed", zap.String("stage", "signin.failed"), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, "signin response read failed")
		log.Error("Orbis train signin failed", zap.String("stage", "signin.failed"), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
		return nil, newProviderError(ErrorInternal, err.Error())
	}
	var parsed signInResponse
	if parseErr := parseResponseBody(body, &parsed); parseErr != nil {
		span.SetStatus(codes.Error, "signin provider error")
		log.Error("Orbis train signin failed",
			zap.String("stage", "signin.failed"),
			zap.Int("error_code", int(parseErr.Code)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, parseErr
	}
	if strings.TrimSpace(parsed.Result.Response.Signature) == "" {
		span.SetStatus(codes.Error, "signin missing signature")
		log.Error("Orbis train signin failed", zap.String("stage", "signin.failed"), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
		return nil, newProviderError(ErrorInternal, "Orbis signin returned empty signature")
	}
	log.Info("Orbis train signin completed",
		zap.String("stage", "signin.complete"),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return &parsed, nil
}

// Search signs in before each search. Its signature is returned for private
// backend storage for future steps of the same booking flow.
func (c *clientImpl) Search(ctx context.Context, request SearchRequest) (*SearchResponse, *ProviderError) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "provider", "orbis/train", "Search")
	defer span.End()
	start := time.Now()
	log := c.logger.WithTrace(ctx, spanID).With(
		zap.String("provider", "orbis"),
		zap.String("product", "train"),
		zap.String("operation", "search"),
	)
	log.Info("Orbis train search started",
		zap.String("stage", "01.validate"),
		zap.String("direction", request.Direction),
		zap.Int("journey_count", len(request.JourneyList)),
		zap.Int("adult_count", request.Adult),
		zap.Int("infant_count", request.Infant),
	)

	if err := request.Validate(); err != nil {
		span.SetStatus(codes.Error, "invalid search request")
		log.Warn("Orbis train search validation failed", zap.String("stage", "01.validate.failed"))
		return nil, newProviderError(ErrorInvalidInput, err.Error())
	}
	log.Info("Orbis train signin step", zap.String("stage", "02.signin"))
	signin, signInErr := c.signIn(ctx)
	if signInErr != nil {
		span.SetStatus(codes.Error, "train signin failed")
		log.Error("Orbis train search failed during signin",
			zap.String("stage", "02.signin.failed"),
			zap.Int("error_code", int(signInErr.Code)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, signInErr
	}

	headers := map[string]string{
		"Action":    "search",
		"Signature": signin.Result.Response.Signature,
	}
	log.Info("Orbis train provider search step", zap.String("stage", "03.search"))
	timeout := searchTimeout
	resp, sendErr := c.sendRequest(ctx, request, bookingTrainPath, &headers, client.RequestOptions{Timeout: &timeout})
	if sendErr != nil {
		span.SetStatus(codes.Error, "train search request failed")
		log.Error("Orbis train search request failed",
			zap.String("stage", "03.search.failed"),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, newProviderError(ErrorInternal, sendErr.Error())
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		span.SetStatus(codes.Error, "train search response read failed")
		log.Error("Orbis train search response failed",
			zap.String("stage", "03.search.failed"),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, newProviderError(ErrorInternal, readErr.Error())
	}
	var result SearchResponse
	if decodeErr := parseResponseBody(body, &result); decodeErr != nil {
		span.SetStatus(codes.Error, "train search provider error")
		log.Error("Orbis train search provider error",
			zap.String("stage", "03.search.failed"),
			zap.Int("error_code", int(decodeErr.Code)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, decodeErr
	}
	result.Signature = signin.Result.Response.Signature
	journeys := 0
	for _, schedule := range result.Result.Response.Schedules {
		journeys += len(schedule.Journeys)
	}
	log.Info("Orbis train search completed",
		zap.String("stage", "04.complete"),
		zap.Int("schedule_count", len(result.Result.Response.Schedules)),
		zap.Int("result_journey_count", journeys),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return &result, nil
}
