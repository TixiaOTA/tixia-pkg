# Redis Package

This package contains Redis-related functionality for the Tixia backend with a shared base client implementation.

## Structure

- `client/` - Shared base Redis client implementation
  - `config.go` - Base Redis configuration structure
  - `base.go` - Base Redis client with common functionality
- `cache/` - Redis cache implementation
  - `cache.go` - Cache manager with TTL, key generation, and pattern-based operations
  - `redis.go` - Redis connection manager for caching
- `streams/` - Redis Streams implementation for event-driven architecture
  - `streams.go` - Main StreamManager implementation with event handling, circuit breakers, dead letter queues, and health monitoring
  - `constants.go` - Event topic constants for different business operations
  - `helpers.go` - Utility functions for working with events and their metadata

## Base Client Package

The `client/` package provides a shared base Redis client implementation between the cache and streams packages.

### BaseRedisConfig

A shared configuration structure for Redis connections:

```go
type BaseRedisConfig struct {
    URL          string        `json:"url" default:"redis://localhost:6379"`
    Password     string        `json:"password"`
    DB           int           `json:"db" default:"0"`
    PoolSize     int           `json:"pool_size" default:"10"`
    MinIdleConns int           `json:"min_idle_conns" default:"5"`
    MaxRetries   int           `json:"max_retries" default:"10"`
    DialTimeout  time.Duration `json:"dial_timeout" default:"5s"`
    ReadTimeout  time.Duration `json:"read_timeout" default:"3s"`
    WriteTimeout time.Duration `json:"write_timeout" default:"3s"`
}
```

### BaseRedisClient

A base client that provides common Redis functionality:

```go
type BaseRedisClient struct {
    client *redis.Client
    logger *zap.Logger
    config BaseRedisConfig
}
```

#### Methods

- `NewBaseRedisClient(config BaseRedisConfig, logger *zap.Logger) (*BaseRedisClient, error)` - Creates a new base Redis client
- `GetClient() *redis.Client` - Returns the underlying Redis client
- `HealthCheck(ctx context.Context) error` - Performs a health check on the Redis connection
- `Close() error` - Closes the Redis connection
- `GetConnectionInfo(ctx context.Context) (map[string]interface{}, error)` - Returns connection information
- `GetStats(ctx context.Context) (map[string]interface{}, error)` - Returns basic Redis statistics

## Environment Configuration

The Redis configuration supports separate Redis servers for cache and streams operations. This provides better workload isolation and flexibility for future scaling.

### Environment Variables

#### Cache Redis Variables

- `REDIS_CACHE_URL` - Cache Redis server URL
- `REDIS_CACHE_PASSWORD` - Cache Redis password
- `REDIS_CACHE_DB` - Cache Redis database number
- `REDIS_CACHE_POOL_SIZE` - Cache connection pool size
- `REDIS_CACHE_MIN_IDLE_CONNS` - Cache minimum idle connections
- `REDIS_CACHE_MAX_RETRIES` - Cache maximum retry attempts

#### Streams Redis Variables

- `REDIS_STREAMS_URL` - Streams Redis server URL
- `REDIS_STREAMS_PASSWORD` - Streams Redis password
- `REDIS_STREAMS_DB` - Streams Redis database number
- `REDIS_STREAMS_POOL_SIZE` - Streams connection pool size
- `REDIS_STREAMS_MIN_IDLE_CONNS` - Streams minimum idle connections
- `REDIS_STREAMS_MAX_RETRIES` - Streams maximum retry attempts

## Usage

### Cache Package

The cache package uses the base client for Redis operations:

```go
import "github.com/TixiaOTA/tixia-pkg/redis/cache"
import "github.com/TixiaOTA/tixia-pkg/redis/client"

// Create cache manager with embedded base config
cacheConfig := cache.CacheConfig{
    BaseRedisConfig: client.BaseRedisConfig{
        URL:          "redis://localhost:6379",
        PoolSize:     10,
        MinIdleConns: 5,
        MaxRetries:   10,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    },
    DefaultTTL: 15 * time.Minute,
    MaxTTL:     24 * time.Hour,
    KeyPrefix:  "tixia:",
}

cacheManager, err := cache.NewCacheManager(cacheConfig, logger)
```

### Streams Package

The streams package provides a robust Redis Streams implementation with the following features:

- **Event Publishing**: Publish events to Redis streams with automatic retry and circuit breaker
- **Event Consumption**: Consume events with consumer groups and automatic reconnection
- **Dead Letter Queue**: Failed events are stored in DLQ for retry
- **Health Monitoring**: Built-in health checks and consumer monitoring
- **Stream Trimming**: Automatic stream trimming to maintain TTL-like behavior
- **Circuit Breaker**: Automatic circuit breaker to prevent cascading failures

#### Example Usage

```go
import "github.com/TixiaOTA/tixia-pkg/redis/streams"
import "github.com/TixiaOTA/tixia-pkg/redis/client"

// Create stream manager with base config
config := streams.StreamConfig{
    BaseRedisConfig: client.BaseRedisConfig{
        URL:          "redis://localhost:6379",
        PoolSize:     10,
        MinIdleConns: 5,
        MaxRetries:   10,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    },
    StreamPrefix:  "tixia",
    ConsumerGroup: "tixia-group",
    ConsumerName:  "tixia-consumer",
}

streamManager, err := streams.NewStreamManager(config, logger)
if err != nil {
    log.Fatal(err)
}
defer streamManager.Close()

// Register event handler
streamManager.RegisterHandler(streams.TopicFlightSearchResults, func(ctx context.Context, event streams.Event) error {
    // Handle the event
    return nil
})

// Publish event
event := streams.Event{
    Type: streams.TopicFlightSearchRequested,
    Data: map[string]interface{}{
        "search_id": "123",
        "from": "JKT",
        "to": "CGK",
    },
    Metadata: map[string]interface{}{
        "request_id": "req_123",
        "user_id": "user_456",
    },
}

err = streamManager.PublishEvent(ctx, event)
if err != nil {
    log.Error(err)
}

// Start consumer
err = streamManager.StartConsumer(ctx, "tixia:flight.search.results")
if err != nil {
    log.Error(err)
}
```

#### Helper Functions

The package includes utility functions for working with events:

```go
// Extract metadata from events
if requestID, ok := streams.GetEventMetadata(event, "request_id"); ok {
    fmt.Printf("Request ID: %s\n", requestID)
}

// Extract nested data from events
if flightData, ok := streams.GetEventData(event, "flight_info"); ok {
    fmt.Printf("Flight data: %+v\n", flightData)
}
```

## Future Plans

- Add comprehensive unit tests for the base client
- Add performance monitoring and metrics collection
- Enhance connection pool management
- Add circuit breaker pattern to base client
- Add configuration validation
- Additional Redis utilities as needed
