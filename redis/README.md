# Redis Package

A comprehensive Redis package for the Tixia backend with shared base client implementation, caching, and event streaming capabilities.

## Overview

The Redis package enables you to:
- Use a shared base Redis client with common functionality
- Implement Redis-based caching with TTL, key generation, and pattern-based operations
- Build event-driven architectures with Redis Streams
- Handle event publishing and consumption with consumer groups
- Implement circuit breakers and dead letter queues for streams
- Monitor health and connection status
- Use separate Redis instances for cache and streams operations

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/redis
```

## Configuration

### Configuration

The Redis package supports separate configurations for cache and streams operations, providing better workload isolation.

#### Cache Configuration

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/redis/cache"
    "github.com/TixiaOTA/tixia-pkg/redis/client"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

cacheConfig := cache.CacheConfig{
    BaseRedisConfig: client.BaseRedisConfig{
        URL:          "redis://localhost:6379",
        Password:     "",
        DB:           0,
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

#### Streams Configuration

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/redis/streams"
    "github.com/TixiaOTA/tixia-pkg/redis/client"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

streamConfig := streams.StreamConfig{
    BaseRedisConfig: client.BaseRedisConfig{
        URL:          "redis://localhost:6379",
        Password:     "",
        DB:           0,
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
    ServiceName:   "tixia-pkg",
    ReadBlock:     30 * time.Second,
    ReadCount:     50,
}

streamManager, err := streams.NewStreamManager(streamConfig, logger)
```

### Configuration Fields

#### BaseRedisConfig

Shared configuration for all Redis operations:

- **URL** (string, required): Redis server URL (e.g., `"redis://localhost:6379"`)
- **Password** (string, optional): Redis password
- **DB** (int, optional): Redis database number. Defaults to 0
- **PoolSize** (int, optional): Connection pool size. Defaults to 10
- **MinIdleConns** (int, optional): Minimum idle connections. Defaults to 5
- **MaxRetries** (int, optional): Maximum retry attempts. Defaults to 10
- **DialTimeout** (time.Duration, optional): Connection timeout. Defaults to 5s
- **ReadTimeout** (time.Duration, optional): Read timeout. Defaults to 3s
- **WriteTimeout** (time.Duration, optional): Write timeout. Defaults to 3s

#### CacheConfig

Cache-specific configuration:

- **BaseRedisConfig** (BaseRedisConfig, required): Base Redis configuration
- **DefaultTTL** (time.Duration, required): Default TTL for cached values
- **MaxTTL** (time.Duration, required): Maximum allowed TTL
- **KeyPrefix** (string, required): Prefix for all cache keys
- **EnableMetrics** (bool, optional): Enable metrics collection. Defaults to false

#### StreamConfig

Streams-specific configuration:

- **BaseRedisConfig** (BaseRedisConfig, required): Base Redis configuration
- **StreamPrefix** (string, optional): Prefix for stream names. Defaults to `"tixia"`
- **ConsumerGroup** (string, optional): Consumer group name. Defaults to `"tixia-group"`
- **ConsumerName** (string, optional): Consumer name. Defaults to `"tixia-consumer"`
- **ServiceName** (string, optional): Service name. Defaults to `"tixia-pkg"`
- **ReadBlock** (time.Duration, optional): Block time for reading streams. Defaults to 30s
- **ReadCount** (int64, optional): Number of entries to read per batch. Defaults to 50
- **ReadErrorBaseDelay** (time.Duration, optional): Base delay for read errors. Defaults to 200ms
- **ReadErrorMaxDelay** (time.Duration, optional): Maximum delay for read errors. Defaults to 5s
- **PublishRetries** (int, optional): Number of publish retries. Defaults to 5
- **PublishBaseDelay** (time.Duration, optional): Base delay for publish retries. Defaults to 50ms
- **PublishMaxDelay** (time.Duration, optional): Maximum delay for publish retries. Defaults to 5s
- **StreamTrimmingEnabled** (bool, optional): Enable automatic stream trimming. Defaults to false
- **MaxStreamLength** (int64, optional): Maximum entries per stream. Defaults to 1000
- **StreamRetentionPeriod** (time.Duration, optional): How long to keep entries. Defaults to 24h

## Usage

### Cache Package

The cache package provides Redis-based caching with TTL, key generation, and pattern-based operations:

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/redis/cache"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

ctx := context.Background()

// Set a value with default TTL
err := cacheManager.Set(ctx, "user:123", userData, nil)

// Set a value with custom TTL
err := cacheManager.Set(ctx, "user:123", userData, &cache.CacheOptions{
    TTL: 1 * time.Hour,
})

// Get a value
var userData User
err := cacheManager.Get(ctx, "user:123", &userData)

// Check if key exists
exists, err := cacheManager.Exists(ctx, "user:123")

// Delete a key
err := cacheManager.Delete(ctx, "user:123")

// Delete by pattern
deleted, err := cacheManager.DeleteByPattern(ctx, "user:*")

// Get TTL
ttl, err := cacheManager.TTL(ctx, "user:123")

// Renew TTL
err := cacheManager.RenewTTL(ctx, "user:123", 1*time.Hour)
```

### Streams Package

The streams package provides Redis Streams implementation for event-driven architecture:

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/redis/streams"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

ctx := context.Background()

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

err := streamManager.PublishEvent(ctx, event)

// Start consumer
err := streamManager.StartConsumer(ctx, "tixia:flight.search.results")
```

## API Reference

### Cache Package

#### NewCacheManager

Creates a new cache manager with the specified configuration.

```go
func NewCacheManager(cacheConfig CacheConfig, logger *logger.Logger) (*CacheManager, error)
```

**Parameters:**

- `cacheConfig` (CacheConfig, required): Cache configuration
- `logger` (*logger.Logger, required): Logger instance

**Return Values:**

- `*CacheManager`: Cache manager instance
- `error`: Error if configuration is invalid or connection fails

#### Methods on CacheManager

##### Set

Sets a value in the cache with optional TTL.

```go
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, opts *CacheOptions) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Cache key
- `value` (interface{}, required): Value to cache (must be JSON serializable)
- `opts` (*CacheOptions, optional): Cache options (TTL, prefix, compression)

**Return Values:**

- `error`: Error if operation fails

##### SetNX

Sets a value only if the key does not already exist. This is an atomic operation and is safe to use for locks, one-time tokens, and replay protection.

```go
func (cm *CacheManager) SetNX(ctx context.Context, key string, value interface{}, opts *CacheOptions) (bool, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Cache key
- `value` (interface{}, required): Value to cache (must be JSON serializable)
- `opts` (*CacheOptions, optional): Cache options (TTL, prefix, compression)

**Return Values:**

- `bool`: `true` when this call created the key, `false` when it already existed
- `error`: Error if operation fails

##### Get

Gets a value from the cache.

```go
func (cm *CacheManager) Get(ctx context.Context, key string, dest interface{}) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Cache key
- `dest` (interface{}, required): Pointer to destination variable

**Return Values:**

- `error`: Error if key not found or operation fails

##### Delete

Deletes a key from the cache.

```go
func (cm *CacheManager) Delete(ctx context.Context, key string) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Cache key to delete

**Return Values:**

- `error`: Error if operation fails

### Streams Package

#### NewStreamManager

Creates a new stream manager with the specified configuration.

```go
func NewStreamManager(config StreamConfig, logger *logger.Logger) (*StreamManager, error)
```

**Parameters:**

- `config` (StreamConfig, required): Stream configuration
- `logger` (*logger.Logger, required): Logger instance

**Return Values:**

- `*StreamManager`: Stream manager instance
- `error`: Error if configuration is invalid or connection fails

#### Methods on StreamManager

##### RegisterHandler

Registers an event handler for a specific event type.

```go
func (sm *StreamManager) RegisterHandler(eventType string, handler EventHandler)
```

**Parameters:**

- `eventType` (string, required): Event type to handle
- `handler` (EventHandler, required): Handler function

##### PublishEvent

Publishes an event to a Redis stream.

```go
func (sm *StreamManager) PublishEvent(ctx context.Context, event Event) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `event` (Event, required): Event to publish

**Return Values:**

- `error`: Error if publish fails

##### StartConsumer

Starts consuming events from a stream.

```go
func (sm *StreamManager) StartConsumer(ctx context.Context, streamName string) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `streamName` (string, required): Stream name to consume from

**Return Values:**

- `error`: Error if consumer fails to start

### Base Client Package

#### NewBaseRedisClient

Creates a new base Redis client.

```go
func NewBaseRedisClient(config BaseRedisConfig, logger *zap.Logger) (*BaseRedisClient, error)
```

**Parameters:**

- `config` (BaseRedisConfig, required): Base Redis configuration
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*BaseRedisClient`: Base Redis client instance
- `error`: Error if connection fails

#### Methods on BaseRedisClient

##### HealthCheck

Performs a health check on the Redis connection.

```go
func (b *BaseRedisClient) HealthCheck(ctx context.Context) error
```

**Parameters:**

- `ctx` (context.Context, required): Context

**Return Values:**

- `error`: Error if health check fails

##### Close

Closes the Redis connection.

```go
func (b *BaseRedisClient) Close() error
```

**Return Values:**

- `error`: Error if close fails

## Error Handling

The package handles errors gracefully:

- **Connection errors**: Returns error during initialization, allowing retry logic
- **Redis unavailable**: Cache operations return errors, streams use circuit breaker pattern
- **Invalid configuration**: Returns error during initialization
- **Serialization errors**: Returns error when marshaling/unmarshaling fails
- **Stream consumption errors**: Uses exponential backoff and retry logic
- **Publish errors**: Retries with exponential backoff up to configured limit

### Common Error Scenarios

- **Redis connection timeout**: Check network connectivity and Redis server status
- **Invalid Redis URL**: Verify URL format (e.g., `redis://localhost:6379`)
- **Authentication failure**: Verify password and Redis ACL settings
- **Key not found**: Normal for cache operations, check with `Exists()` first
- **Stream consumer group errors**: Ensure consumer group exists or use `XGROUP CREATE`

## Data Models

### BaseRedisConfig

Shared configuration structure for Redis connections:

```go
type BaseRedisConfig struct {
    URL          string        `json:"url"`
    Password     string        `json:"password"`
    DB           int           `json:"db"`
    PoolSize     int           `json:"pool_size"`
    MinIdleConns int           `json:"min_idle_conns"`
    MaxRetries   int           `json:"max_retries"`
    DialTimeout  time.Duration `json:"dial_timeout"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}
```

### CacheConfig

Cache manager configuration:

```go
type CacheConfig struct {
    BaseRedisConfig client.BaseRedisConfig
    DefaultTTL      time.Duration
    MaxTTL          time.Duration
    KeyPrefix       string
    EnableMetrics   bool
}
```

### CacheOptions

Options for individual cache operations:

```go
type CacheOptions struct {
    TTL         time.Duration
    KeyPrefix   string
    Compression bool
}
```

### StreamConfig

Stream manager configuration:

```go
type StreamConfig struct {
    BaseRedisConfig       client.BaseRedisConfig
    StreamPrefix          string
    ConsumerGroup         string
    ConsumerName          string
    ServiceName           string
    ReadBlock             time.Duration
    ReadCount             int64
    ReadErrorBaseDelay    time.Duration
    ReadErrorMaxDelay     time.Duration
    PublishRetries        int
    PublishBaseDelay      time.Duration
    PublishMaxDelay       time.Duration
    StreamTrimmingEnabled bool
    MaxStreamLength       int64
    StreamRetentionPeriod  time.Duration
}
```

### Event

Represents a generic event structure:

```go
type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Source    string                 `json:"source"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
```

## Examples

### Complete Cache Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    "log"
    
    "github.com/TixiaOTA/tixia-pkg/redis/cache"
    "github.com/TixiaOTA/tixia-pkg/redis/client"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    logger := logger.NewLogger()
    
    cacheConfig := cache.CacheConfig{
        BaseRedisConfig: client.BaseRedisConfig{
            URL:          "redis://localhost:6379",
            PoolSize:     10,
            MinIdleConns: 5,
            DialTimeout:  5 * time.Second,
        },
        DefaultTTL: 15 * time.Minute,
        MaxTTL:     24 * time.Hour,
        KeyPrefix:  "tixia:",
    }
    
    cacheManager, err := cache.NewCacheManager(cacheConfig, logger)
    if err != nil {
        log.Fatal(err)
    }
    defer cacheManager.Close()
    
    ctx := context.Background()
    
    // Set user data
    user := User{ID: 123, Name: "John Doe", Email: "john@example.com"}
    err = cacheManager.Set(ctx, "user:123", user, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get user data
    var cachedUser User
    err = cacheManager.Get(ctx, "user:123", &cachedUser)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Cached user: %+v\n", cachedUser)
}
```

### Complete Streams Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/TixiaOTA/tixia-pkg/redis/streams"
    "github.com/TixiaOTA/tixia-pkg/redis/client"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

func main() {
    logger := logger.NewLogger()
    
    streamConfig := streams.StreamConfig{
        BaseRedisConfig: client.BaseRedisConfig{
            URL:          "redis://localhost:6379",
            PoolSize:     10,
            DialTimeout:  5 * time.Second,
        },
        StreamPrefix:  "tixia",
        ConsumerGroup: "tixia-group",
        ConsumerName:  "tixia-consumer",
    }
    
    streamManager, err := streams.NewStreamManager(streamConfig, logger)
    if err != nil {
        log.Fatal(err)
    }
    defer streamManager.Close()
    
    // Register handler
    streamManager.RegisterHandler(streams.TopicFlightSearchResults, func(ctx context.Context, event streams.Event) error {
        fmt.Printf("Received event: %+v\n", event)
        return nil
    })
    
    ctx := context.Background()
    
    // Publish event
    event := streams.Event{
        Type: streams.TopicFlightSearchRequested,
        Data: map[string]interface{}{
            "search_id": "123",
            "from": "JKT",
            "to": "CGK",
        },
    }
    
    err = streamManager.PublishEvent(ctx, event)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start consumer
    err = streamManager.StartConsumer(ctx, "tixia:flight.search.results")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Environment Configuration

The Redis configuration supports separate Redis servers for cache and streams operations via environment variables:

### Cache Redis Variables

- `REDIS_CACHE_URL` - Cache Redis server URL
- `REDIS_CACHE_PASSWORD` - Cache Redis password
- `REDIS_CACHE_DB` - Cache Redis database number
- `REDIS_CACHE_POOL_SIZE` - Cache connection pool size
- `REDIS_CACHE_MIN_IDLE_CONNS` - Cache minimum idle connections
- `REDIS_CACHE_MAX_RETRIES` - Cache maximum retry attempts

### Streams Redis Variables

- `REDIS_STREAMS_URL` - Streams Redis server URL
- `REDIS_STREAMS_PASSWORD` - Streams Redis password
- `REDIS_STREAMS_DB` - Streams Redis database number
- `REDIS_STREAMS_POOL_SIZE` - Streams connection pool size
- `REDIS_STREAMS_MIN_IDLE_CONNS` - Streams minimum idle connections
- `REDIS_STREAMS_MAX_RETRIES` - Streams maximum retry attempts

## License

This package is part of the TixiaOTA project.
