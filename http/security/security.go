package security

import (
	"github.com/TixiaOTA/tixia-backend-master/pkg/http/security/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// SecurityManager manages HTTP security middleware
type SecurityManager struct {
	config *Config
	logger *zap.Logger
	redis  *redis.Client
}

// SecurityBuilder provides a fluent interface for building security middleware
type SecurityBuilder struct {
	config *Config
	logger *zap.Logger
	redis  *redis.Client
}

// NewSecurityBuilder creates a new SecurityBuilder with default configuration
func NewSecurityBuilder() *SecurityBuilder {
	return &SecurityBuilder{
		config: DefaultConfig(),
	}
}

// WithConfig sets the complete security configuration
func (b *SecurityBuilder) WithConfig(config *Config) *SecurityBuilder {
	b.config = config
	return b
}

// WithRateLimit configures rate limiting
func (b *SecurityBuilder) WithRateLimit(config RateLimitConfig) *SecurityBuilder {
	b.config.RateLimit = &config
	return b
}

// WithSecurityHeaders configures security headers
func (b *SecurityBuilder) WithSecurityHeaders(config HeadersConfig) *SecurityBuilder {
	b.config.Headers = &config
	return b
}

// WithCORS configures CORS
func (b *SecurityBuilder) WithCORS(config CORSConfig) *SecurityBuilder {
	b.config.CORS = &config
	return b
}

// WithIPFilter configures IP filtering
func (b *SecurityBuilder) WithIPFilter(config IPFilterConfig) *SecurityBuilder {
	b.config.IPFilter = &config
	return b
}

// WithValidation configures request validation
func (b *SecurityBuilder) WithValidation(config ValidationConfig) *SecurityBuilder {
	b.config.Validation = &config
	return b
}

// WithLogger sets the logger instance
func (b *SecurityBuilder) WithLogger(logger *zap.Logger) *SecurityBuilder {
	b.logger = logger
	return b
}

// WithRedis sets the Redis client for rate limiting
func (b *SecurityBuilder) WithRedis(client *redis.Client) *SecurityBuilder {
	b.redis = client
	return b
}

// EnableRateLimit enables rate limiting with default settings
func (b *SecurityBuilder) EnableRateLimit() *SecurityBuilder {
	b.config.RateLimit.Enabled = true
	return b
}

// EnableSecurityHeaders enables security headers with default settings
func (b *SecurityBuilder) EnableSecurityHeaders() *SecurityBuilder {
	b.config.Headers.Enabled = true
	return b
}

// EnableCORS enables CORS with default settings
func (b *SecurityBuilder) EnableCORS() *SecurityBuilder {
	b.config.CORS.Enabled = true
	return b
}

// EnableIPFilter enables IP filtering with default settings
func (b *SecurityBuilder) EnableIPFilter() *SecurityBuilder {
	b.config.IPFilter.Enabled = true
	return b
}

// EnableValidation enables request validation with default settings
func (b *SecurityBuilder) EnableValidation() *SecurityBuilder {
	b.config.Validation.Enabled = true
	return b
}

// OnlyRateLimit disables all other middleware and enables only rate limiting
func (b *SecurityBuilder) OnlyRateLimit(config RateLimitConfig) *SecurityBuilder {
	b.config = DefaultConfig()
	b.config.RateLimit = &config
	b.config.RateLimit.Enabled = true
	return b
}

// OnlyHeaders disables all other middleware and enables only security headers
func (b *SecurityBuilder) OnlyHeaders(config HeadersConfig) *SecurityBuilder {
	b.config = DefaultConfig()
	b.config.Headers = &config
	b.config.Headers.Enabled = true
	return b
}

// OnlyIPFilter disables all other middleware and enables only IP filtering
func (b *SecurityBuilder) OnlyIPFilter(config IPFilterConfig) *SecurityBuilder {
	b.config = DefaultConfig()
	b.config.IPFilter = &config
	b.config.IPFilter.Enabled = true
	return b
}

// OnlyValidation disables all other middleware and enables only validation
func (b *SecurityBuilder) OnlyValidation(config ValidationConfig) *SecurityBuilder {
	b.config = DefaultConfig()
	b.config.Validation = &config
	b.config.Validation.Enabled = true
	return b
}

// ForAPI configures security for API endpoints (rate limiting + validation)
func (b *SecurityBuilder) ForAPI() *SecurityBuilder {
	b.config.RateLimit.Enabled = true
	b.config.Headers.Enabled = false
	b.config.IPFilter.Enabled = false
	b.config.Validation.Enabled = true
	return b
}

// ForWeb configures security for web applications (headers + validation)
func (b *SecurityBuilder) ForWeb() *SecurityBuilder {
	b.config.RateLimit.Enabled = false
	b.config.Headers.Enabled = true
	b.config.IPFilter.Enabled = false
	b.config.Validation.Enabled = true
	return b
}

// ForPublic configures minimal security for public endpoints
func (b *SecurityBuilder) ForPublic() *SecurityBuilder {
	b.config.RateLimit.Enabled = false
	b.config.Headers.Enabled = true
	b.config.IPFilter.Enabled = false
	b.config.Validation.Enabled = true
	return b
}

// ForAdmin configures strict security for admin endpoints
func (b *SecurityBuilder) ForAdmin() *SecurityBuilder {
	b.config.RateLimit.Enabled = true
	b.config.Headers.Enabled = true
	b.config.IPFilter.Enabled = true
	b.config.Validation.Enabled = true

	// Stricter rate limiting for admin
	b.config.RateLimit.Rate = 50
	b.config.RateLimit.Burst = 5

	return b
}

// Build creates a SecurityManager with the configured settings
func (b *SecurityBuilder) Build() *SecurityManager {
	return &SecurityManager{
		config: b.config,
		logger: b.logger,
		redis:  b.redis,
	}
}

// Middleware returns all enabled middleware in the correct order
func (sm *SecurityManager) Middleware() []fiber.Handler {
	var handlers []fiber.Handler

	// CORS should be first (before rate limiting)
	if sm.config.CORS != nil && sm.config.CORS.Enabled {
		handlers = append(handlers, sm.CORSMiddleware())
	}

	// Order matters: IP filter first, then rate limit, then headers, finally validation
	if sm.config.IPFilter != nil && sm.config.IPFilter.Enabled {
		handlers = append(handlers, sm.IPFilterMiddleware())
	}

	if sm.config.RateLimit != nil && sm.config.RateLimit.Enabled {
		handlers = append(handlers, sm.RateLimitMiddleware())
	}

	if sm.config.Headers != nil && sm.config.Headers.Enabled {
		handlers = append(handlers, sm.HeadersMiddleware())
	}

	if sm.config.Validation != nil && sm.config.Validation.Enabled {
		handlers = append(handlers, sm.ValidationMiddleware())
	}

	return handlers
}

// Individual middleware getters
func (sm *SecurityManager) RateLimitMiddleware() fiber.Handler {
	if sm.config.RateLimit == nil || !sm.config.RateLimit.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	rateLimiter := middleware.NewRateLimiter(sm.redis, sm.logger, middleware.RateLimitConfig{
		Rate:    sm.config.RateLimit.Rate,
		Window:  sm.config.RateLimit.Window,
		Burst:   sm.config.RateLimit.Burst,
		KeyFunc: sm.config.RateLimit.KeyFunc,
		Message: sm.config.RateLimit.Message,
	})

	return rateLimiter.Middleware()
}

func (sm *SecurityManager) CORSMiddleware() fiber.Handler {
	if sm.config.CORS == nil || !sm.config.CORS.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	cors := middleware.NewCORS(sm.logger, middleware.CORSConfig{
		Enabled:             sm.config.CORS.Enabled,
		AllowOrigins:        sm.config.CORS.AllowOrigins,
		AllowMethods:        sm.config.CORS.AllowMethods,
		AllowHeaders:        sm.config.CORS.AllowHeaders,
		AllowCredentials:    sm.config.CORS.AllowCredentials,
		ExposeHeaders:       sm.config.CORS.ExposeHeaders,
		MaxAge:              sm.config.CORS.MaxAge,
		AllowPrivateNetwork: sm.config.CORS.AllowPrivateNetwork,
		SecurityHeaders:     sm.config.CORS.SecurityHeaders,
	})

	return cors.Middleware()
}

func (sm *SecurityManager) HeadersMiddleware() fiber.Handler {
	if sm.config.Headers == nil || !sm.config.Headers.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	headers := middleware.NewSecurityHeaders(sm.logger, middleware.HeadersConfig{
		CSP:                  sm.config.Headers.CSP,
		HSTS:                 sm.config.Headers.HSTS,
		HSTSMaxAge:           sm.config.Headers.HSTSMaxAge,
		NoSniff:              sm.config.Headers.NoSniff,
		FrameOptions:         sm.config.Headers.FrameOptions,
		XSSProtection:        sm.config.Headers.XSSProtection,
		ReferrerPolicy:       sm.config.Headers.ReferrerPolicy,
		PermittedCrossDomain: sm.config.Headers.PermittedCrossDomain,
		CustomHeaders:        sm.config.Headers.CustomHeaders,
	})

	return headers.Middleware()
}

func (sm *SecurityManager) IPFilterMiddleware() fiber.Handler {
	// TODO: Implement IP filtering middleware
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func (sm *SecurityManager) ValidationMiddleware() fiber.Handler {
	// TODO: Implement request validation middleware
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// Convenience factory functions

// ForAPI creates security manager optimized for API endpoints
func ForAPI(redis *redis.Client, logger *zap.Logger) *SecurityManager {
	return NewSecurityBuilder().
		ForAPI().
		EnableCORS().
		WithRedis(redis).
		WithLogger(logger).
		Build()
}

// ForWeb creates security manager optimized for web applications
func ForWeb(redis *redis.Client, logger *zap.Logger) *SecurityManager {
	return NewSecurityBuilder().
		ForWeb().
		EnableCORS().
		WithRedis(redis).
		WithLogger(logger).
		Build()
}

// ForPublic creates security manager for public endpoints
func ForPublic() *SecurityManager {
	return NewSecurityBuilder().
		ForPublic().
		EnableCORS().
		Build()
}

// ForAdmin creates security manager with strict settings for admin endpoints
func ForAdmin(redis *redis.Client, logger *zap.Logger) *SecurityManager {
	return NewSecurityBuilder().
		ForAdmin().
		WithRedis(redis).
		WithLogger(logger).
		Build()
}

// RateLimitOnly creates a middleware that only does rate limiting
func RateLimitOnly(config RateLimitConfig, redis *redis.Client, logger *zap.Logger) fiber.Handler {
	return NewSecurityBuilder().
		OnlyRateLimit(config).
		WithRedis(redis).
		WithLogger(logger).
		Build().
		RateLimitMiddleware()
}

// CORSOnly creates a middleware that only handles CORS
func CORSOnly(config CORSConfig, logger *zap.Logger) fiber.Handler {
	return NewSecurityBuilder().
		WithCORS(config).
		WithLogger(logger).
		Build().
		CORSMiddleware()
}
