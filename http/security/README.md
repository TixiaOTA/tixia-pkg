# HTTP Security Package

A simple and flexible HTTP security middleware package for Go Fiber applications, designed with a builder pattern for easy configuration and use.

## Overview

The HTTP Security package enables you to:
- Apply Redis-backed rate limiting with sliding window algorithm
- Add OWASP recommended security headers
- Configure CORS with security considerations
- Filter requests by IP address (whitelist/blacklist)
- Validate and sanitize requests (coming soon)
- Use flexible builder pattern or preset configurations
- Apply security middleware globally or per-route

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/http/security
```

## Configuration

### Configuration

Create a security manager using the builder pattern:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/http/security"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
)

// Using builder pattern
security := security.NewSecurityBuilder().
    WithRateLimit(security.RateLimitConfig{
        Enabled: true,
        Rate:    100,              // 100 requests
        Window:  time.Minute,      // per minute
        Burst:   10,               // burst of 10
        KeyFunc: "ip",             // rate limit by IP
    }).
    WithSecurityHeaders(security.HeadersConfig{Enabled: true}).
    WithRedis(redisClient).
    WithLogger(logger).
    Build()

// Or use convenience factories
apiSecurity := security.ForAPI(redisClient, logger)
publicSecurity := security.ForPublic()
adminSecurity := security.ForAdmin(redisClient, logger)
```

### Configuration Fields

#### RateLimitConfig

- **Enabled** (bool, required): Enable/disable rate limiting
- **Rate** (int, optional): Requests per window. Defaults to 100
- **Window** (time.Duration, optional): Time window duration. Defaults to 1 minute
- **Burst** (int, optional): Burst capacity. Defaults to 10
- **KeyFunc** (string, optional): Rate limiting key function. Options: `"ip"`, `"user"`, `"endpoint"`, `"user_endpoint"`. Defaults to `"ip"`
- **Message** (string, optional): Custom rate limit exceeded message. Defaults to "Rate limit exceeded. Please try again later."

#### HeadersConfig

- **Enabled** (bool, required): Enable/disable security headers
- **CSP** (string, optional): Content Security Policy. Defaults to `"default-src 'self'"`
- **HSTS** (bool, optional): Enable HTTP Strict Transport Security. Defaults to true
- **HSTSMaxAge** (int, optional): HSTS max age in seconds. Defaults to 31536000 (1 year)
- **NoSniff** (bool, optional): Enable X-Content-Type-Options: nosniff. Defaults to true
- **FrameOptions** (string, optional): X-Frame-Options value. Defaults to `"DENY"`
- **XSSProtection** (bool, optional): Enable X-XSS-Protection. Defaults to true
- **ReferrerPolicy** (string, optional): Referrer-Policy value. Defaults to `"strict-origin-when-cross-origin"`
- **PermittedCrossDomain** (bool, optional): Enable X-Permitted-Cross-Domain-Policies. Defaults to false
- **CustomHeaders** (map[string]string, optional): Custom security headers

#### CORSConfig

- **Enabled** (bool, required): Enable/disable CORS
- **AllowOrigins** ([]string, optional): Allowed origins. Defaults to localhost origins
- **AllowMethods** ([]string, optional): Allowed HTTP methods. Defaults to common methods
- **AllowHeaders** ([]string, optional): Allowed headers. Defaults to common headers
- **AllowCredentials** (bool, optional): Allow credentials. Defaults to true
- **ExposeHeaders** ([]string, optional): Headers to expose
- **MaxAge** (int, optional): Max age for preflight requests in seconds. Defaults to 86400 (24 hours)
- **AllowPrivateNetwork** (bool, optional): Allow private network access. Defaults to false
- **SecurityHeaders** (bool, optional): Add security headers to CORS responses. Defaults to true

#### IPFilterConfig

- **Enabled** (bool, required): Enable/disable IP filtering
- **Whitelist** ([]string, optional): CIDR blocks or IPs to allow
- **Blacklist** ([]string, optional): CIDR blocks or IPs to deny
- **Mode** (string, optional): `"whitelist"` or `"blacklist"`. Defaults to `"blacklist"`
- **Message** (string, optional): Custom blocked message. Defaults to "Access denied from your IP address."

#### ValidationConfig

- **Enabled** (bool, required): Enable/disable request validation
- **MaxBodySize** (int64, optional): Maximum request body size in bytes. Defaults to 10MB
- **AllowedMethods** ([]string, optional): Allowed HTTP methods
- **RequiredHeaders** ([]string, optional): Headers that must be present
- **BlockedPaths** ([]string, optional): Paths to block (regex patterns)
- **AllowedContentTypes** ([]string, optional): Allowed content types
- **ValidateJSON** (bool, optional): Validate JSON payloads. Defaults to true
- **ValidateXML** (bool, optional): Validate XML payloads. Defaults to false

## Usage

### Basic Rate Limiting Only

```go
import "github.com/TixiaOTA/tixia-pkg/http/security"

// Enable only rate limiting
security := security.NewSecurityBuilder().
    WithRateLimit(security.RateLimitConfig{
        Enabled: true,
        Rate:    100,              // 100 requests
        Window:  time.Minute,      // per minute
        Burst:   10,               // burst of 10
        KeyFunc: "ip",             // rate limit by IP
    }).
    WithSecurityHeaders(security.HeadersConfig{Enabled: false}).
    WithIPFilter(security.IPFilterConfig{Enabled: false}).
    WithValidation(security.ValidationConfig{Enabled: false}).
    WithRedis(redisClient).
    WithLogger(logger).
    Build()

app.Use(security.Middleware()...)
```

### Using Convenience Factories

```go
// For API endpoints
apiSecurity := security.ForAPI(redisClient, logger)
app.Use("/api", apiSecurity.Middleware()...)

// For public endpoints
publicSecurity := security.ForPublic()
app.Use("/public", publicSecurity.Middleware()...)

// For admin endpoints (strict security)
adminSecurity := security.ForAdmin(redisClient, logger)
app.Use("/admin", adminSecurity.Middleware()...)
```

### Individual Middleware

```go
security := security.NewSecurityBuilder().
    WithRateLimit(security.RateLimitConfig{Enabled: true, Rate: 100}).
    WithRedis(redisClient).
    Build()

// Use only specific middleware
app.Use(security.RateLimitMiddleware())
app.Use(security.HeadersMiddleware())
```

### Route-Specific Security

Apply different security configurations to different route groups:

```go
func setupSecurityRoutes(app *fiber.App, redisClient *redis.Client, logger *zap.Logger) {
    // Public routes - no security
    public := app.Group("/public")
    public.Get("/health", healthHandler)

    // API routes - only rate limiting
    api := app.Group("/api")
    apiSecurity := security.NewSecurityBuilder().
        WithRateLimit(security.RateLimitConfig{
            Enabled: true,
            Rate:    100,
            Window:  time.Minute,
            KeyFunc: "ip",
        }).
        WithSecurityHeaders(security.HeadersConfig{Enabled: false}).
        WithIPFilter(security.IPFilterConfig{Enabled: false}).
        WithValidation(security.ValidationConfig{Enabled: false}).
        WithRedis(redisClient).
        WithLogger(logger).
        Build()

    for _, middleware := range apiSecurity.Middleware() {
        api.Use(middleware)
    }
    api.Get("/users", getUsersHandler)
    api.Post("/bookings", createBookingHandler)

    // Admin routes - strict security
    admin := app.Group("/admin")
    adminSecurity := security.NewSecurityBuilder().
        WithRateLimit(security.RateLimitConfig{
            Enabled: true,
            Rate:    50,              // Stricter rate limit
            Window:  time.Minute,
            KeyFunc: "user",          // Rate limit by user
        }).
        WithSecurityHeaders(security.HeadersConfig{
            Enabled: true,
            HSTS:    true,
            CSP:     "default-src 'self'",
        }).
        WithRedis(redisClient).
        WithLogger(logger).
        Build()

    for _, middleware := range adminSecurity.Middleware() {
        admin.Use(middleware)
    }
    admin.Get("/dashboard", adminDashboardHandler)
}
```

### Endpoint-Specific Rate Limiting

Different rate limits for different endpoints:

```go
// Strict rate limiting for login endpoint
loginRateLimit := security.RateLimitOnly(
    security.RateLimitConfig{
        Enabled: true,
        Rate:    5,                // Only 5 attempts per minute
        Window:  time.Minute,
        KeyFunc: "ip",
        Message: "Too many login attempts. Please try again later.",
    },
    redisClient,
    logger,
)
app.Post("/auth/login", loginRateLimit, loginHandler)

// Lenient rate limiting for search endpoint
searchRateLimit := security.RateLimitOnly(
    security.RateLimitConfig{
        Enabled: true,
        Rate:    1000,             // 1000 requests per minute
        Window:  time.Minute,
        KeyFunc: "ip",
    },
    redisClient,
    logger,
)
app.Get("/api/search", searchRateLimit, searchHandler)
```

## API Reference

### Public Functions

#### NewSecurityBuilder

Creates a new SecurityBuilder with default configuration.

```go
func NewSecurityBuilder() *SecurityBuilder
```

**Return Values:**

- `*SecurityBuilder`: A new security builder instance

#### ForAPI

Creates security manager optimized for API endpoints (rate limiting + validation + CORS).

```go
func ForAPI(redis *redis.Client, logger *zap.Logger) *SecurityManager
```

**Parameters:**

- `redis` (*redis.Client, required): Redis client for rate limiting
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*SecurityManager`: Configured security manager for API endpoints

#### ForWeb

Creates security manager optimized for web applications (security headers + validation + CORS).

```go
func ForWeb(redis *redis.Client, logger *zap.Logger) *SecurityManager
```

**Parameters:**

- `redis` (*redis.Client, required): Redis client
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*SecurityManager`: Configured security manager for web applications

#### ForPublic

Creates security manager for public endpoints (minimal security).

```go
func ForPublic() *SecurityManager
```

**Return Values:**

- `*SecurityManager`: Configured security manager for public endpoints

#### ForAdmin

Creates security manager with strict settings for admin endpoints.

```go
func ForAdmin(redis *redis.Client, logger *zap.Logger) *SecurityManager
```

**Parameters:**

- `redis` (*redis.Client, required): Redis client for rate limiting
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*SecurityManager`: Configured security manager with strict settings

#### RateLimitOnly

Creates a middleware that only does rate limiting.

```go
func RateLimitOnly(config RateLimitConfig, redis *redis.Client, logger *zap.Logger) fiber.Handler
```

**Parameters:**

- `config` (RateLimitConfig, required): Rate limiting configuration
- `redis` (*redis.Client, required): Redis client
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `fiber.Handler`: Rate limiting middleware handler

### Public Types

#### SecurityManager

Manages HTTP security middleware.

```go
type SecurityManager struct {
    config *Config
    logger *zap.Logger
    redis  *redis.Client
}
```

#### Methods on SecurityManager

##### Middleware

Returns all enabled middleware in the correct order.

```go
func (sm *SecurityManager) Middleware() []fiber.Handler
```

**Return Values:**

- `[]fiber.Handler`: Slice of middleware handlers

##### RateLimitMiddleware

Returns rate limiting middleware if enabled.

```go
func (sm *SecurityManager) RateLimitMiddleware() fiber.Handler
```

**Return Values:**

- `fiber.Handler`: Rate limiting middleware, or nil if disabled

##### HeadersMiddleware

Returns security headers middleware if enabled.

```go
func (sm *SecurityManager) HeadersMiddleware() fiber.Handler
```

**Return Values:**

- `fiber.Handler`: Security headers middleware, or nil if disabled

#### SecurityBuilder

Provides a fluent interface for building security middleware.

```go
type SecurityBuilder struct {
    config *Config
    logger *zap.Logger
    redis  *redis.Client
}
```

#### Methods on SecurityBuilder

##### WithRateLimit

Configures rate limiting.

```go
func (b *SecurityBuilder) WithRateLimit(config RateLimitConfig) *SecurityBuilder
```

**Parameters:**

- `config` (RateLimitConfig, required): Rate limiting configuration

**Return Values:**

- `*SecurityBuilder`: Builder instance for method chaining

##### WithSecurityHeaders

Configures security headers.

```go
func (b *SecurityBuilder) WithSecurityHeaders(config HeadersConfig) *SecurityBuilder
```

**Parameters:**

- `config` (HeadersConfig, required): Security headers configuration

**Return Values:**

- `*SecurityBuilder`: Builder instance for method chaining

##### WithRedis

Sets the Redis client for rate limiting.

```go
func (b *SecurityBuilder) WithRedis(client *redis.Client) *SecurityBuilder
```

**Parameters:**

- `client` (*redis.Client, required): Redis client instance

**Return Values:**

- `*SecurityBuilder`: Builder instance for method chaining

##### WithLogger

Sets the logger instance.

```go
func (b *SecurityBuilder) WithLogger(logger *zap.Logger) *SecurityBuilder
```

**Parameters:**

- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*SecurityBuilder`: Builder instance for method chaining

##### Build

Creates a SecurityManager with the configured settings.

```go
func (b *SecurityBuilder) Build() *SecurityManager
```

**Return Values:**

- `*SecurityManager`: Configured security manager

## Error Handling

The package gracefully handles errors:

- **Redis connection issues**: Logs warning and allows requests to continue
- **Invalid configuration**: Uses sensible defaults
- **Rate limit errors**: Logs error and allows requests (does not block on Redis errors)
- **IP filter errors**: Logs warning and allows requests if Redis is unavailable

## Data Models

### Config

Holds all security middleware configuration.

```go
type Config struct {
    RateLimit  *RateLimitConfig  `json:"rate_limit"`
    Headers    *HeadersConfig    `json:"headers"`
    CORS       *CORSConfig       `json:"cors"`
    IPFilter   *IPFilterConfig   `json:"ip_filter"`
    Validation *ValidationConfig `json:"validation"`
}
```

### RateLimitConfig

Configures rate limiting middleware.

```go
type RateLimitConfig struct {
    Enabled bool          `json:"enabled"`
    Rate    int           `json:"rate"`     // requests per window
    Window  time.Duration `json:"window"`   // time window duration
    Burst   int           `json:"burst"`    // burst capacity
    KeyFunc string        `json:"key_func"` // "ip", "user", "endpoint"
    Message string        `json:"message"`  // custom rate limit exceeded message
}
```

### HeadersConfig

Configures security headers middleware.

```go
type HeadersConfig struct {
    Enabled              bool              `json:"enabled"`
    CSP                  string            `json:"csp"`
    HSTS                 bool              `json:"hsts"`
    HSTSMaxAge           int               `json:"hsts_max_age"`
    NoSniff              bool              `json:"no_sniff"`
    FrameOptions         string            `json:"frame_options"`
    XSSProtection        bool              `json:"xss_protection"`
    ReferrerPolicy       string            `json:"referrer_policy"`
    PermittedCrossDomain bool              `json:"permitted_cross_domain"`
    CustomHeaders        map[string]string `json:"custom_headers"`
}
```

### CORSConfig

Configures CORS middleware with security considerations.

```go
type CORSConfig struct {
    Enabled             bool     `json:"enabled"`
    AllowOrigins        []string `json:"allow_origins"`
    AllowMethods        []string `json:"allow_methods"`
    AllowHeaders        []string `json:"allow_headers"`
    AllowCredentials    bool     `json:"allow_credentials"`
    ExposeHeaders       []string `json:"expose_headers"`
    MaxAge              int      `json:"max_age"`
    AllowPrivateNetwork bool     `json:"allow_private_network"`
    SecurityHeaders     bool     `json:"security_headers"`
}
```

### IPFilterConfig

Configures IP filtering middleware.

```go
type IPFilterConfig struct {
    Enabled   bool     `json:"enabled"`
    Whitelist []string `json:"whitelist"` // CIDR blocks or IPs
    Blacklist []string `json:"blacklist"` // CIDR blocks or IPs
    Mode      string   `json:"mode"`      // "whitelist" or "blacklist"
    Message   string   `json:"message"`   // custom blocked message
}
```

### ValidationConfig

Configures request validation middleware.

```go
type ValidationConfig struct {
    Enabled             bool     `json:"enabled"`
    MaxBodySize         int64    `json:"max_body_size"`
    AllowedMethods      []string `json:"allowed_methods"`
    RequiredHeaders     []string `json:"required_headers"`
    BlockedPaths        []string `json:"blocked_paths"`
    AllowedContentTypes []string `json:"allowed_content_types"`
    ValidateJSON        bool     `json:"validate_json"`
    ValidateXML         bool     `json:"validate_xml"`
}
```

## Examples

### Complete Example

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/TixiaOTA/tixia-pkg/http/security"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "time"
)

func main() {
    app := fiber.New()
    logger, _ := zap.NewProduction()
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create security manager for API endpoints
    apiSecurity := security.ForAPI(redisClient, logger)
    
    // Apply to API routes
    api := app.Group("/api/v1")
    for _, middleware := range apiSecurity.Middleware() {
        api.Use(middleware)
    }
    
    api.Get("/users", getUsersHandler)
    api.Post("/bookings", createBookingHandler)

    // Strict rate limiting for sensitive endpoints
    paymentLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    3,
            Window:  time.Minute,
            KeyFunc: "user",
            Message: "Too many payment attempts. Please contact support.",
        },
        redisClient,
        logger,
    )
    api.Post("/payments", paymentLimit, processPaymentHandler)

    app.Listen(":3000")
}
```

## Rate Limiting Headers

When rate limiting is enabled, the following headers are set:

- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets
- `Retry-After`: Seconds to wait when limit exceeded

## Security Headers Applied

When security headers are enabled, the following headers are set:

- `Content-Security-Policy`: Configurable CSP policy
- `Strict-Transport-Security`: HSTS with configurable max-age
- `X-Content-Type-Options`: Set to "nosniff"
- `X-Frame-Options`: Configurable frame options
- `X-XSS-Protection`: XSS protection header
- `Referrer-Policy`: Configurable referrer policy
- Custom headers as configured

## License

This package is part of the TixiaOTA project.
