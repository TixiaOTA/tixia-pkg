package streams

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"github.com/TixiaOTA/tixia-pkg/redis/client"
	"github.com/TixiaOTA/tixia-pkg/redis/redisutil"
	"go.opentelemetry.io/otel/codes"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Event represents a generic event structure
type Event struct {
	ID        string                 `json:"id"`                 // unique identifier for the event
	Type      string                 `json:"type"`               // type of the event (e.g. "flight.search.requested")
	Source    string                 `json:"source"`             // source of the event (e.g. "tixia-pkg")
	Timestamp time.Time              `json:"timestamp"`          // timestamp of the event
	Data      map[string]interface{} `json:"data"`               // data of the event (e.g. "flight search request")
	Metadata  map[string]interface{} `json:"metadata,omitempty"` // metadata of the event (e.g. "request_id", "search_id")
}

// StreamConfig holds Redis Streams configuration
type StreamConfig struct {
	// Base Redis configuration
	BaseRedisConfig client.BaseRedisConfig `envconfig:"REDIS_"`

	// Stream-specific settings
	StreamPrefix  string `envconfig:"STREAM_PREFIX" default:"tixia"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP" default:"tixia-group"`
	ConsumerName  string `envconfig:"CONSUMER_NAME" default:"tixia-consumer"`
	ServiceName   string `envconfig:"SERVICE_NAME" default:"tixia-pkg"`

	// Read settings
	ReadBlock time.Duration `envconfig:"REDIS_READ_BLOCK" default:"30s"`
	ReadCount int64         `envconfig:"REDIS_READ_COUNT" default:"50"`

	// Read error backoff settings
	ReadErrorBaseDelay time.Duration `envconfig:"REDIS_READ_BASE_DELAY" default:"200ms"`
	ReadErrorMaxDelay  time.Duration `envconfig:"REDIS_READ_MAX_DELAY" default:"5s"`

	// Reconnection settings
	PublishRetries   int           `envconfig:"REDIS_PUBLISH_RETRIES" default:"5"`
	PublishBaseDelay time.Duration `envconfig:"REDIS_PUBLISH_BASE_DELAY" default:"50ms"`
	PublishMaxDelay  time.Duration `envconfig:"REDIS_PUBLISH_MAX_DELAY" default:"5s"`

	// Stream trimming settings (TTL-like behavior)
	StreamTrimmingEnabled bool          `envconfig:"REDIS_STREAM_TRIMMING_ENABLED" default:"false"`
	MaxStreamLength       int64         `envconfig:"REDIS_MAX_STREAM_LENGTH" default:"1000"`      // Maximum number of entries per stream
	StreamRetentionPeriod time.Duration `envconfig:"REDIS_STREAM_RETENTION_PERIOD" default:"24h"` // How long to keep entries
}

// EventHandler is a function that processes events
type EventHandler func(ctx context.Context, event Event) error

// EventCorrelation represents a correlation between a request and its response
type EventCorrelation struct {
	RequestID  string
	SearchID   string
	ResponseCh chan Event
	Timeout    time.Time
}

// StreamManager manages Redis Streams operations
type StreamManager struct {
	client       *redis.Client
	config       StreamConfig
	logger       *logger.Logger
	handlers     map[string]EventHandler
	correlations map[string]*EventCorrelation // requestID -> correlation
	mu           sync.RWMutex
}

// NewStreamManager creates a new Redis Streams manager
func NewStreamManager(config StreamConfig, logger *logger.Logger) (*StreamManager, error) {
	baseClient, err := client.NewBaseRedisClient(config.BaseRedisConfig, logger.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create base Redis client: %w", err)
	}

	sm := &StreamManager{
		client:       baseClient.GetClient(),
		config:       config,
		logger:       logger,
		handlers:     make(map[string]EventHandler),
		correlations: make(map[string]*EventCorrelation),
	}

	return sm, nil
}

// TrimStreamWithCustomLimits trims a Redis stream with custom limits
func (sm *StreamManager) TrimStreamWithCustomLimits(ctx context.Context, streamName string, maxEntries int64, retentionPeriod time.Duration) error {
	if maxEntries <= 0 && retentionPeriod <= 0 {
		return fmt.Errorf("at least one of maxEntries or retentionPeriod must be greater than 0")
	}

	// Method 1: Trim by length (keep only the most recent N entries)
	if maxEntries > 0 {
		result := sm.client.XTrimMaxLen(ctx, streamName, maxEntries)
		if result.Err() != nil {
			return fmt.Errorf("failed to trim stream by length: %w", result.Err())
		}
		sm.logger.Debug("Stream trimmed by length with custom limit",
			zap.String("stream", streamName),
			zap.Int64("max_length", maxEntries),
			zap.Int64("trimmed_count", result.Val()),
		)
	}

	// Method 2: Trim by time (keep entries newer than retention period)
	if retentionPeriod > 0 {
		cutoffTime := time.Now().UTC().Add(-retentionPeriod)
		cutoffID := fmt.Sprintf("%d-0", cutoffTime.UnixMilli())

		result := sm.client.XTrimMinID(ctx, streamName, cutoffID)
		if result.Err() != nil {
			return fmt.Errorf("failed to trim stream by time: %w", result.Err())
		}
		sm.logger.Debug("Stream trimmed by time with custom limit",
			zap.String("stream", streamName),
			zap.String("cutoff_id", cutoffID),
			zap.Duration("retention_period", retentionPeriod),
			zap.Int64("trimmed_count", result.Val()),
		)
	}

	return nil
}

// Close gracefully closes the Redis connection and cleans up resources
func (sm *StreamManager) Close() error {
	sm.logger.Info("Shutting down stream manager...")

	// Close all correlation channels to prevent goroutine leaks
	sm.mu.Lock()
	for requestID, correlation := range sm.correlations {
		close(correlation.ResponseCh)
		sm.logger.Debug("Closed correlation channel", zap.String("request_id", requestID))
	}
	// Clear correlations map
	sm.correlations = make(map[string]*EventCorrelation)
	sm.mu.Unlock()

	// Close Redis client connection
	if sm.client != nil {
		if err := sm.client.Close(); err != nil {
			sm.logger.Error("Failed to close Redis client", zap.Error(err))
			return fmt.Errorf("failed to close Redis client: %w", err)
		}
		sm.logger.Info("Redis client connection closed successfully")
	}

	sm.logger.Info("Stream manager shutdown completed")
	return nil
}

// RegisterHandler registers an event handler for a specific event type
func (sm *StreamManager) RegisterHandler(eventType string, handler EventHandler) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.handlers[eventType] = handler
	sm.logger.Info("Registered event handler", zap.String("event_type", eventType))
}

// RegisterCorrelation registers a correlation between a request and its expected response
func (sm *StreamManager) RegisterCorrelation(requestID string, searchID string, timeout time.Duration) chan Event {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	responseCh := make(chan Event, 1)
	correlation := &EventCorrelation{
		RequestID:  requestID,
		SearchID:   searchID,
		ResponseCh: responseCh,
		Timeout:    time.Now().UTC().Add(timeout),
	}

	sm.correlations[requestID] = correlation
	sm.logger.Debug("Registered event correlation",
		zap.String("request_id", requestID),
		zap.String("search_id", searchID),
		zap.Duration("timeout", timeout),
	)

	return responseCh
}

// UnregisterCorrelation removes a correlation and closes its channel
func (sm *StreamManager) UnregisterCorrelation(requestID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if correlation, exists := sm.correlations[requestID]; exists {
		close(correlation.ResponseCh)
		delete(sm.correlations, requestID)
		sm.logger.Debug("Unregistered event correlation",
			zap.String("request_id", requestID),
		)
	}
}

// GetCorrelation retrieves a correlation by request ID
func (sm *StreamManager) GetCorrelation(requestID string) (*EventCorrelation, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	correlation, exists := sm.correlations[requestID]
	return correlation, exists
}

// PublishEvent publishes an event to a Redis stream
func (sm *StreamManager) PublishEvent(ctx context.Context, event Event) error {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "handler", "PublishEvent", "PublishEvent")
	defer span.End()

	streamName := fmt.Sprintf("%s:%s", sm.config.StreamPrefix, event.Type)

	if event.ID == "" {
		event.ID = GenerateEventID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.Source == "" {
		event.Source = sm.config.ServiceName
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Reconnection configuration for publishing
	maxRetries := sm.config.PublishRetries
	baseDelay := sm.config.PublishBaseDelay
	maxDelay := sm.config.PublishMaxDelay
	currentDelay := baseDelay
	retryCount := 0
	traceMap := otelpkg.InjectTraceContext(ctx)
	traceJSON, _ := json.Marshal(traceMap)

	for retryCount <= maxRetries {
		result := sm.client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event":        string(eventData),
				"traceContext": traceJSON,
			},
		})

		if result.Err() == nil {
			// Success - log and return
			sm.logger.Info("Event published",
				zap.String("stream", streamName),
				zap.String("event_id", event.ID),
				zap.String("event_type", event.Type),
			)

			return nil
		}

		// Handle error with retry logic
		if retryCount >= maxRetries {
			span.SetStatus(codes.Error, "Max retries exceeded for publishing event")
			sm.logger.WithTrace(ctx, currentSpanId).Error("Max retries exceeded for publishing event",
				zap.String("stream", streamName),
				zap.String("event_id", event.ID),
				zap.String("event_type", event.Type),
				zap.Error(result.Err()),
			)
			return fmt.Errorf("failed to publish event after %d retries: %w", maxRetries, result.Err())
		}

		sm.logger.WithTrace(ctx, currentSpanId).Warn("Failed to publish event, retrying",
			zap.String("stream", streamName),
			zap.String("event_id", event.ID),
			zap.String("event_type", event.Type),
			zap.Error(result.Err()),
			zap.Int("retry_count", retryCount),
			zap.Duration("delay", currentDelay),
		)

		// Wait before retrying
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(currentDelay):
			// Continue to retry
		}

		// Exponential backoff with jitter
		retryCount++
		currentDelay = time.Duration(float64(currentDelay) * 1.5)
		if currentDelay > maxDelay {
			currentDelay = maxDelay
		}

		// Add jitter to prevent thundering herd
		jitter := time.Duration(float64(currentDelay) * 0.1 * (0.5 + 0.5*float64(time.Now().UTC().UnixNano()%100)/100))
		currentDelay += jitter

		// Test connection before retrying
		if err := sm.testConnection(ctx); err != nil {
			sm.logger.WithTrace(ctx, currentSpanId).Error("Connection test failed, will retry",
				zap.String("stream", streamName),
				zap.Error(err),
			)
			continue
		}
	}

	return fmt.Errorf("failed to publish event after retries")
}

func (sm *StreamManager) CreateConsumerGroup(ctx context.Context, streamName string) error {
	groups, err := sm.client.XInfoGroups(ctx, streamName).Result()
	if err != nil {
		_, err = sm.client.XGroupCreateMkStream(ctx, streamName, sm.config.ConsumerGroup, "0").Result()
		if err != nil {
			return fmt.Errorf("failed to create consumer group: %w", err)
		}
		sm.logger.Info("Created consumer group with stream", zap.String("stream", streamName))
		return nil
	}

	for _, group := range groups {
		if group.Name == sm.config.ConsumerGroup {
			sm.logger.Info("Consumer group already exists", zap.String("stream", streamName))
			return nil
		}
	}

	_, err = sm.client.XGroupCreate(ctx, streamName, sm.config.ConsumerGroup, "0").Result()
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	sm.logger.Info("Created consumer group", zap.String("stream", streamName))
	return nil
}

func (sm *StreamManager) StartConsumer(ctx context.Context, streamName string, timeout time.Duration) error {
	if err := sm.CreateConsumerGroup(ctx, streamName); err != nil {
		return err
	}

	sm.logger.Info("Starting consumer", zap.String("stream", streamName))

	currentDelay := sm.config.ReadErrorBaseDelay
	maxDelay := sm.config.ReadErrorMaxDelay

	for {
		select {
		case <-ctx.Done():
			sm.logger.Info("Consumer shutdown requested", zap.String("stream", streamName))
			return ctx.Err()
		default:
			streams, err := sm.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    sm.config.ConsumerGroup,
				Consumer: sm.config.ConsumerName,
				Streams:  []string{streamName, ">"},
				Count:    sm.config.ReadCount,
				Block:    sm.config.ReadBlock,
			}).Result()

			if err != nil {
				switch {
				case err == redis.Nil:
					continue
				case errors.Is(err, context.Canceled):
					sm.logger.Info("Consumer stopping due to context cancellation", zap.String("stream", streamName))
					return ctx.Err()
				case errors.Is(err, context.DeadlineExceeded):
					// Normal when Block elapses with no messages
					sm.logger.Debug("Read block deadline reached; continuing", zap.String("stream", streamName))
					continue
				default:
					sm.logger.Warn("Error reading from stream; backing off",
						zap.Error(err),
						zap.String("stream", streamName),
						zap.Duration("backoff", currentDelay),
					)
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(currentDelay):
					}
					// Exponential backoff with jitter
					currentDelay = time.Duration(float64(currentDelay) * 1.5)
					if currentDelay > maxDelay {
						currentDelay = maxDelay
					}
					jitter := time.Duration(float64(currentDelay) * 0.1 * (0.5 + 0.5*float64(time.Now().UTC().UnixNano()%100)/100))
					currentDelay += jitter
					continue
				}
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					// Check if context is still valid before processing
					if ctx.Err() != nil {
						sm.logger.Info("Context canceled, stopping message processing", zap.String("stream", streamName))
						return ctx.Err()
					}

					// consumer.go
					rawTraceCtx := message.Values["traceContext"]
					traceMap := redisutil.NormalizeTraceMap(rawTraceCtx)
					ctxWithTrace := otelpkg.ExtractTraceContext(ctx, traceMap)

					if err := sm.processMessage(ctxWithTrace, message); err != nil {
						sm.logger.Error("Failed to process message",
							zap.String("message_id", message.ID),
							zap.Error(err),
						)
					}
				}
			}
		}
	}
}

// StartConsumerForStreams starts a single consumer that reads from multiple streams in one loop
func (sm *StreamManager) StartConsumerForStreams(ctx context.Context, streamNames []string, timeout time.Duration) error {
	for _, s := range streamNames {
		if err := sm.CreateConsumerGroup(ctx, s); err != nil {
			return err
		}
	}

	// Build XREADGROUP Streams argument for go-redis:
	// first N entries are stream names, followed by N IDs (usually ">")
	streamsArg := make([]string, 0, len(streamNames)*2)
	// append names
	streamsArg = append(streamsArg, streamNames...)
	// append IDs
	for range streamNames {
		streamsArg = append(streamsArg, ">")
	}

	sm.logger.Info("Starting multi-stream consumer", zap.Strings("streams", streamNames))

	currentDelay := sm.config.ReadErrorBaseDelay
	maxDelay := sm.config.ReadErrorMaxDelay

	for {
		select {
		case <-ctx.Done():
			sm.logger.Info("Multi-stream consumer shutdown requested")
			return ctx.Err()
		default:
			results, err := sm.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    sm.config.ConsumerGroup,
				Consumer: sm.config.ConsumerName,
				Streams:  streamsArg,
				Count:    sm.config.ReadCount,
				Block:    sm.config.ReadBlock,
			}).Result()

			if err != nil {
				switch {
				case err == redis.Nil:
					continue
				case errors.Is(err, context.Canceled):
					sm.logger.Info("Multi-stream consumer stopping due to context cancellation")
					return ctx.Err()
				case errors.Is(err, context.DeadlineExceeded):
					sm.logger.Debug("Multi-stream read block deadline reached; continuing")
					continue
				default:
					sm.logger.Warn("Error reading from multiple streams; backing off", zap.Error(err), zap.Duration("backoff", currentDelay))
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(currentDelay):
					}
					currentDelay = time.Duration(float64(currentDelay) * 1.5)
					if currentDelay > maxDelay {
						currentDelay = maxDelay
					}
					jitter := time.Duration(float64(currentDelay) * 0.1 * (0.5 + 0.5*float64(time.Now().UTC().UnixNano()%100)/100))
					currentDelay += jitter
					continue
				}
			}

			for _, stream := range results {
				for _, message := range stream.Messages {
					if ctx.Err() != nil {
						return ctx.Err()
					}

					rawTraceCtx := message.Values["traceContext"]
					traceMap := redisutil.NormalizeTraceMap(rawTraceCtx)
					ctxWithTrace := otelpkg.ExtractTraceContext(ctx, traceMap)

					if err := sm.processMessage(ctxWithTrace, message); err != nil {
						sm.logger.Error("Failed to process message", zap.String("message_id", message.ID), zap.Error(err))
					}
				}
			}
		}
	}
}

// testConnection performs a simple ping test to verify Redis connectivity
func (sm *StreamManager) testConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return sm.client.Ping(ctx).Err()
}

func (sm *StreamManager) processMessage(ctx context.Context, message redis.XMessage) error {
	// Early check for context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	eventData, ok := message.Values["event"].(string)
	if !ok {
		return fmt.Errorf("invalid event data format")
	}

	var event Event
	if err := json.Unmarshal([]byte(eventData), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	streamName := fmt.Sprintf("%s:%s", sm.config.StreamPrefix, event.Type)

	// Check if this event has a correlation (request-response pattern)
	if event.Metadata != nil {
		if requestID, ok := event.Metadata["request_id"].(string); ok {
			if correlation, exists := sm.GetCorrelation(requestID); exists {
				// Check if the correlation is still valid (not timed out)
				if time.Now().UTC().Before(correlation.Timeout) {
					// Route the event to the correlation channel
					select {
					case correlation.ResponseCh <- event:
						sm.logger.Debug("Event routed to correlation channel",
							zap.String("request_id", requestID),
						)
					default:
						sm.logger.Warn("Correlation channel is full or closed",
							zap.String("request_id", requestID),
							zap.String("event_type", event.Type),
						)
					}
				} else {
					sm.logger.Warn("Correlation timed out, removing",
						zap.String("request_id", requestID),
						zap.String("event_type", event.Type),
					)
					sm.UnregisterCorrelation(requestID)
				}
				// Always acknowledge the message with a fresh context to avoid cancellation issues
				ackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := sm.client.XAck(ackCtx, streamName, sm.config.ConsumerGroup, message.ID).Result()
				cancel()
				if err != nil {
					return fmt.Errorf("failed to acknowledge message: %w", err)
				}
				return nil
			}
		}
	}

	// Process with registered handlers if no correlation found
	sm.mu.RLock()
	handler, exists := sm.handlers[event.Type]
	sm.mu.RUnlock()

	if !exists {
		sm.logger.Warn("No handler found for event type", zap.String("event_type", event.Type))
		// Still acknowledge the message even if no handler exists
		ackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := sm.client.XAck(ackCtx, streamName, sm.config.ConsumerGroup, message.ID).Result()
		cancel()
		return err
	}

	// consumer.go
	rawTraceCtx := message.Values["traceContext"]
	traceMap := redisutil.NormalizeTraceMap(rawTraceCtx)
	ctxWithTrace := otelpkg.ExtractTraceContext(ctx, traceMap)
	// Create a fresh context for handler execution to prevent cancellation issues
	handlerCtx, cancel := context.WithTimeout(ctxWithTrace, 3*time.Minute)

	handlerErr := handler(handlerCtx, event)
	cancel()

	if handlerErr != nil {
		sm.logger.Error("Handler failed",
			zap.String("event_type", event.Type),
			zap.Error(handlerErr),
		)
	}

	// Always acknowledge the message with a fresh context
	ackCtx, ackCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := sm.client.XAck(ackCtx, streamName, sm.config.ConsumerGroup, message.ID).Result()
	ackCancel()
	if err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	sm.logger.Info("Message processed and acknowledged",
		zap.String("message_id", message.ID),
		zap.String("event_type", event.Type),
	)

	return nil
}

func (sm *StreamManager) WaitForStreamResponseAndUnmarshal(ctx context.Context, timeout time.Duration, requestId, searchId string, target any) error {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "stream manager", "WaitForStreamResponseAndUnmarshal", "WaitForStreamResponseAndUnmarshal")
	defer span.End()

	// Wait for the response from the provider
	responseCh := sm.RegisterCorrelation(requestId, searchId, timeout)
	defer sm.UnregisterCorrelation(requestId)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case response := <-responseCh:
		sm.logger.WithTrace(ctx, currentSpanId).Info("response received", zap.Any("event", response.Data))
		responseBytes, err := json.Marshal(&response.Data)
		if err != nil {
			sm.logger.WithTrace(ctx, currentSpanId).Error("failed to marshal response", zap.Error(err))
			return err
		}

		err = json.Unmarshal(responseBytes, &target)
		if err != nil {
			sm.logger.WithTrace(ctx, currentSpanId).Error("failed to unmarshal response", zap.Error(err))
			return err
		}

		sm.logger.WithTrace(ctx, currentSpanId).Info("unmarshal success", zap.Any("response", response))

	case <-ctx.Done():
		sm.logger.WithTrace(ctx, currentSpanId).Error("timeout waiting for response", zap.String("request_id", requestId))
		return ctx.Err()
	}

	return nil
}
