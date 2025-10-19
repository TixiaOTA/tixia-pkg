package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CORSConfig holds CORS configuration with security considerations
type CORSConfig struct {
	Enabled             bool     // Enable CORS
	AllowOrigins        []string // Allowed origins
	AllowMethods        []string // Allowed HTTP methods
	AllowHeaders        []string // Allowed headers
	AllowCredentials    bool     // Allow credentials
	ExposeHeaders       []string // Headers to expose
	MaxAge              int      // Max age for preflight requests
	AllowPrivateNetwork bool     // Allow private network access
	SecurityHeaders     bool     // Add security headers to CORS responses
}

// CORS implements CORS middleware with security enhancements
type CORS struct {
	logger *zap.Logger
	config CORSConfig
}

// NewCORS creates a new CORS middleware instance
func NewCORS(logger *zap.Logger, config CORSConfig) *CORS {
	return &CORS{
		logger: logger,
		config: config,
	}
}

// Middleware returns the CORS middleware
func (c *CORS) Middleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Skip CORS for non-browser requests
		if !c.isBrowserRequest(ctx) {
			return ctx.Next()
		}

		origin := ctx.Get("Origin")
		if origin == "" {
			return ctx.Next()
		}

		// Check if origin is allowed
		if !c.isOriginAllowed(origin) {
			c.logger.Warn("CORS: Origin not allowed",
				zap.String("origin", origin),
				zap.String("ip", ctx.IP()),
				zap.String("user_agent", ctx.Get("User-Agent")),
			)
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Origin not allowed",
			})
		}

		// Set CORS headers
		ctx.Set("Access-Control-Allow-Origin", origin)

		if c.config.AllowCredentials {
			ctx.Set("Access-Control-Allow-Credentials", "true")
		}

		if len(c.config.AllowMethods) > 0 {
			ctx.Set("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
		}

		if len(c.config.AllowHeaders) > 0 {
			ctx.Set("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
		}

		if len(c.config.ExposeHeaders) > 0 {
			ctx.Set("Access-Control-Expose-Headers", strings.Join(c.config.ExposeHeaders, ", "))
		}

		if c.config.MaxAge > 0 {
			ctx.Set("Access-Control-Max-Age", string(rune(c.config.MaxAge)))
		}

		if c.config.AllowPrivateNetwork {
			ctx.Set("Access-Control-Allow-Private-Network", "true")
		}

		// Add security headers to CORS responses if enabled
		if c.config.SecurityHeaders {
			c.addSecurityHeaders(ctx)
		}

		// Handle preflight requests
		if ctx.Method() == fiber.MethodOptions {
			return ctx.SendStatus(fiber.StatusNoContent)
		}

		return ctx.Next()
	}
}

// isBrowserRequest checks if the request is from a browser
func (c *CORS) isBrowserRequest(ctx *fiber.Ctx) bool {
	userAgent := ctx.Get("User-Agent")
	return strings.Contains(userAgent, "Mozilla") ||
		strings.Contains(userAgent, "Chrome") ||
		strings.Contains(userAgent, "Safari") ||
		strings.Contains(userAgent, "Firefox") ||
		strings.Contains(userAgent, "Edge")
}

// isOriginAllowed checks if the origin is in the allowed list
func (c *CORS) isOriginAllowed(origin string) bool {
	if !c.config.Enabled {
		return false
	}

	// Allow all origins if none specified (not recommended for production)
	if len(c.config.AllowOrigins) == 0 {
		return false
	}

	// Check for wildcard
	if len(c.config.AllowOrigins) == 1 && c.config.AllowOrigins[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedOrigin := range c.config.AllowOrigins {
		if allowedOrigin == origin {
			return true
		}
	}

	return false
}

// addSecurityHeaders adds security headers to CORS responses
func (c *CORS) addSecurityHeaders(ctx *fiber.Ctx) {
	// Basic security headers for CORS responses
	ctx.Set("X-Content-Type-Options", "nosniff")
	ctx.Set("X-Frame-Options", "DENY")
	ctx.Set("X-XSS-Protection", "1; mode=block")
	ctx.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	// Remove server information
	ctx.Set("X-Powered-By", "")
	ctx.Set("Server", "")
}
