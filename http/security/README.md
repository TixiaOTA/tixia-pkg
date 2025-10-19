# HTTP Security Package

A simple and flexible HTTP security middleware package for Go Fiber applications, designed with a builder pattern for easy configuration and use.

## Features

- **Rate Limiting**: Redis-backed rate limiting with sliding window algorithm
- **Security Headers**: OWASP recommended security headers
- **Request Validation**: Input validation and sanitization (coming soon)
- **IP Filtering**: IP whitelist/blacklist support (coming soon)
- **Flexible Configuration**: Builder pattern with preset configurations
- **Production Ready**: Comprehensive logging and error handling

## Quick Start

### Basic Rate Limiting Only

```go
import "github.com/TixiaOTA/tixia-backend-master/pkg/http/security"

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
```

## Route-Specific Security

### Approach 1: Route Groups (Recommended)

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

### Approach 2: Individual Middleware per Route

Apply specific middleware to individual routes:

```go
func setupIndividualSecurity(app *fiber.App, redisClient *redis.Client, logger *zap.Logger) {
    
    // Create security manager
    security := security.NewSecurityBuilder().
        WithRateLimit(security.RateLimitConfig{
            Enabled: true,
            Rate:    100,
            Window:  time.Minute,
        }).
        WithSecurityHeaders(security.HeadersConfig{
            Enabled: true,
        }).
        WithRedis(redisClient).
        WithLogger(logger).
        Build()
    
    // Route with only rate limiting
    app.Get("/api/public", 
        security.RateLimitMiddleware(),
        publicAPIHandler,
    )
    
    // Route with only security headers
    app.Get("/webapp/profile", 
        security.HeadersMiddleware(),
        profileHandler,
    )
    
    // Route with both rate limiting and headers
    app.Post("/api/sensitive", 
        security.RateLimitMiddleware(),
        security.HeadersMiddleware(),
        sensitiveAPIHandler,
    )
    
    // Route with no security
    app.Get("/health", healthHandler)
}
```

### Approach 3: Endpoint-Specific Rate Limiting

Different rate limits for different endpoints:

```go
func setupEndpointSpecificRateLimits(app *fiber.App, redisClient *redis.Client, logger *zap.Logger) {
    
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
    
    // Very strict for payment processing
    paymentRateLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    5,                // Only 5 payment attempts per minute
            Window:  time.Minute,
            KeyFunc: "user",           // Per user, not IP
            Message: "Too many payment attempts. Please contact support.",
        },
        redisClient,
        logger,
    )
    app.Post("/api/payments", paymentRateLimit, paymentHandler)
}
```

### Approach 4: Tixia-Specific Integration

Integration with your existing Tixia route structure:

```go
func SetupTixiaRoutesWithSecurity(app *fiber.App, h *handlers.Handlers, redisClient *redis.Client, logger *zap.Logger) {
    
    // Health routes - no security
    app.Get("/", handlers.GetServiceInfo)
    app.Get("/health", handlers.HealthCheck)
    app.Get("/health/db", handlers.DBHealthCheck)
    
    // API v1 routes with standard rate limiting
    apiV1 := app.Group("/api/v1")
    apiSecurity := security.ForAPI(redisClient, logger)
    for _, middleware := range apiSecurity.Middleware() {
        apiV1.Use(middleware)
    }
    
    // Setup your existing routes with security
    setupAirportRoutes(apiV1, h)
    setupCityRoutes(apiV1, h)
    setupCountryRoutes(apiV1, h)
    setupFlightRoutes(apiV1, h)
    
    // Booking routes with special rate limiting
    bookingRateLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    10,               // Only 10 bookings per minute per user
            Window:  time.Minute,
            KeyFunc: "user",
            Message: "You're creating bookings too quickly. Please wait.",
        },
        redisClient,
        logger,
    )
    apiV1.Post("/bookings", bookingRateLimit, h.CreateBooking)
    apiV1.Get("/bookings", h.GetBookings)
    apiV1.Get("/bookings/:id", h.GetBookingByID)
    
    // Payment routes with very strict rate limiting
    paymentRateLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    3,                // Only 3 payment attempts per minute
            Window:  time.Minute,
            KeyFunc: "user",
            Message: "Too many payment attempts. Please contact support if you need assistance.",
        },
        redisClient,
        logger,
    )
    apiV1.Post("/payments", paymentRateLimit, h.ProcessPayment)
    apiV1.Get("/payments/:id", h.GetPaymentStatus)
    
    // Public documentation routes with minimal security
    publicSecurity := security.ForPublic()
    docsGroup := app.Group("/docs")
    for _, middleware := range publicSecurity.Middleware() {
        docsGroup.Use(middleware)
    }
    docsGroup.Get("/swagger/*", swagger.HandlerDefault)
}
```

### Key Recommendations

1. **Use Route Groups** for broad categorizations (public, api, admin)
2. **Use Individual Middleware** for specific routes that need different security
3. **Use Endpoint-Specific Rate Limiting** for sensitive operations like:
   - Login/Authentication: 5-10 attempts per minute
   - Booking Creation: 10-20 per minute per user
   - Payment Processing: 3-5 per minute per user
   - Search/Read Operations: 100-1000 per minute
   - Admin Operations: 50-100 per minute

4. **Security Order**: Apply security middleware before your business logic handlers
5. **Rate Limiting Keys**:
   - Use `"ip"` for anonymous endpoints
   - Use `"user"` for authenticated endpoints
   - Use `"endpoint"` for per-endpoint limits
   - Use `"user_endpoint"` for per-user per-endpoint limits

## Configuration

### Rate Limiting

```go
RateLimitConfig{
    Enabled: true,
    Rate:    100,                    // requests per window
    Window:  time.Minute,            // time window
    Burst:   10,                     // burst capacity
    KeyFunc: "ip",                   // "ip", "user", "endpoint", "user_endpoint"
    Message: "Rate limit exceeded",  // custom message
}
```

### Security Headers

```go
HeadersConfig{
    Enabled:              true,
    CSP:                  "default-src 'self'",
    HSTS:                 true,
    HSTSMaxAge:           31536000,
    NoSniff:              true,
    FrameOptions:         "DENY",
    XSSProtection:        true,
    ReferrerPolicy:       "strict-origin-when-cross-origin",
    PermittedCrossDomain: false,
    CustomHeaders:        map[string]string{
        "X-API-Version": "v1.0",
    },
}
```

## Preset Configurations

- **ForAPI()**: Rate limiting + validation (ideal for REST APIs)
- **ForWeb()**: Security headers + validation (ideal for web apps)
- **ForPublic()**: Minimal security for public endpoints
- **ForAdmin()**: Strict security for admin endpoints

## Integration with Existing Tixia Backend

### Option 1: Global Security (Simple)

Add to your server setup in `internal/api/server.go`:

```go
func NewServer(config ServerConfig, logger *logger.Logger, ...) *Server {
    app := fiber.New(fiber.Config{...})
    
    // Existing middleware
    app.Use(recover.New(...))
    app.Use(requestid.New(...))
    
    // Add global security middleware
    if config.Environment == "production" {
        securityManager := security.ForAPI(redisClient, logger.Logger)
        for _, middleware := range securityManager.Middleware() {
            app.Use(middleware)
        }
    }
    
    // Rest of your middleware
    app.Use(cors.New(...))
    app.Use(middleware.LoggerMiddleware(logger))
    
    return &Server{...}
}
```

### Option 2: Route-Specific Security (Recommended)

Modify your route setup in `internal/api/routes/routes.go`:

```go
func SetupRoutes(app *fiber.App, h *handlers.Handlers, dbManager *db.ConnectionManager, 
                streamManager *streams.StreamManager, jwtSecretUser, jwtSecretBackoffice string, 
                logger *logger.Logger, redisClient *redis.Client) {

    // Health and info routes - no security
    app.Get("/", handlers.GetServiceInfo)
    app.Get("/health", handlers.HealthCheck)
    app.Get("/health/db", handlers.DBHealthCheck)

    // API v1 group with standard rate limiting
    apiV1 := app.Group("/api/v1")
    apiSecurity := security.ForAPI(redisClient, logger.Logger)
    for _, middleware := range apiSecurity.Middleware() {
        apiV1.Use(middleware)
    }

    // Standard API routes
    setupAirportRoutes(apiV1, h)
    setupCityRoutes(apiV1, h)
    setupCountryRoutes(apiV1, h)
    setupFlightRoutes(apiV1, h)

    // Booking routes with stricter rate limiting
    setupBookingRoutesWithSecurity(apiV1, h, redisClient, logger.Logger)
    
    // Payment routes with very strict rate limiting  
    setupPaymentRoutesWithSecurity(apiV1, h, redisClient, logger.Logger)

    // Documentation routes with minimal security
    docsGroup := app.Group("/docs")
    publicSecurity := security.ForPublic()
    for _, middleware := range publicSecurity.Middleware() {
        docsGroup.Use(middleware)
    }
    docsGroup.Get("/swagger/*", swagger.HandlerDefault)
}

func setupBookingRoutesWithSecurity(router fiber.Router, h *handlers.Handlers, 
                                   redisClient *redis.Client, logger *zap.Logger) {
    
    // Special rate limit for booking creation
    bookingCreateLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    10,               // 10 bookings per minute per user
            Window:  time.Minute,
            KeyFunc: "user",
            Message: "You're creating bookings too quickly. Please wait.",
        },
        redisClient,
        logger,
    )
    
    router.Post("/bookings", bookingCreateLimit, h.CreateBooking)
    router.Get("/bookings", h.GetBookings)
    router.Get("/bookings/:id", h.GetBookingByID)
    router.Put("/bookings/:id", h.UpdateBooking)
    router.Delete("/bookings/:id", h.DeleteBooking)
}

func setupPaymentRoutesWithSecurity(router fiber.Router, h *handlers.Handlers, 
                                   redisClient *redis.Client, logger *zap.Logger) {
    
    // Very strict rate limit for payments
    paymentLimit := security.RateLimitOnly(
        security.RateLimitConfig{
            Enabled: true,
            Rate:    3,                // Only 3 payment attempts per minute
            Window:  time.Minute,
            KeyFunc: "user",
            Message: "Too many payment attempts. Please contact support.",
        },
        redisClient,
        logger,
    )
    
    router.Post("/payments", paymentLimit, h.ProcessPayment)
    router.Get("/payments/:id", h.GetPaymentStatus)
}
```

## Environment-Based Configuration

```go
func NewSecurityFromEnv(env string, redis *redis.Client, logger *zap.Logger) *security.SecurityManager {
    switch env {
    case "production":
        return security.ForAPI(redis, logger)
    case "development":
        return security.ForPublic()
    default:
        return security.NewSecurityBuilder().Build()
    }
}
```

## Error Handling

The package gracefully handles errors:
- Redis connection issues: Logs warning and allows requests
- Invalid configuration: Uses sensible defaults
- Rate limit errors: Logs error and allows requests

## Rate Limiting Headers

When rate limiting is enabled, the following headers are set:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets
- `Retry-After`: Seconds to wait when limit exceeded

## Security Headers Applied

- `Content-Security-Policy`: Configurable CSP policy
- `Strict-Transport-Security`: HSTS with configurable max-age
- `X-Content-Type-Options`: Set to "nosniff"
- `X-Frame-Options`: Configurable frame options
- `X-XSS-Protection`: XSS protection header
- `Referrer-Policy`: Configurable referrer policy
- Custom headers as configured

## Coming Soon

- **Request Validation**: Body size limits, content type validation, JSON/XML validation
- **IP Filtering**: CIDR-based IP whitelist/blacklist support
- **Advanced Rate Limiting**: Per-endpoint rate limits, burst algorithms
- **Metrics**: Prometheus metrics for monitoring

## License

This package is part of the Tixia backend project.
