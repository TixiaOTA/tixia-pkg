package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// SecurityHeaders implements security headers middleware
type SecurityHeaders struct {
	logger *zap.Logger
	config HeadersConfig
}

// HeadersConfig holds security headers configuration
type HeadersConfig struct {
	CSP                  string            // Content Security Policy
	HSTS                 bool              // HTTP Strict Transport Security
	HSTSMaxAge           int               // HSTS max age in seconds
	NoSniff              bool              // X-Content-Type-Options: nosniff
	FrameOptions         string            // X-Frame-Options
	XSSProtection        bool              // X-XSS-Protection
	ReferrerPolicy       string            // Referrer-Policy
	PermittedCrossDomain bool              // X-Permitted-Cross-Domain-Policies
	CustomHeaders        map[string]string // Custom security headers
}

// NewSecurityHeaders creates a new security headers middleware instance
func NewSecurityHeaders(logger *zap.Logger, config HeadersConfig) *SecurityHeaders {
	return &SecurityHeaders{
		logger: logger,
		config: config,
	}
}

// Middleware returns the security headers middleware
func (sh *SecurityHeaders) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Content Security Policy
		if sh.config.CSP != "" {
			c.Set("Content-Security-Policy", sh.config.CSP)
		}

		// HTTP Strict Transport Security
		if sh.config.HSTS {
			maxAge := sh.config.HSTSMaxAge
			if maxAge <= 0 {
				maxAge = 31536000 // Default: 1 year
			}
			hstsValue := fmt.Sprintf("max-age=%d; includeSubDomains", maxAge)
			c.Set("Strict-Transport-Security", hstsValue)
		}

		// X-Content-Type-Options
		if sh.config.NoSniff {
			c.Set("X-Content-Type-Options", "nosniff")
		}

		// X-Frame-Options
		if sh.config.FrameOptions != "" {
			c.Set("X-Frame-Options", sh.config.FrameOptions)
		}

		// X-XSS-Protection
		if sh.config.XSSProtection {
			c.Set("X-XSS-Protection", "1; mode=block")
		}

		// Referrer-Policy
		if sh.config.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", sh.config.ReferrerPolicy)
		}

		// X-Permitted-Cross-Domain-Policies
		if !sh.config.PermittedCrossDomain {
			c.Set("X-Permitted-Cross-Domain-Policies", "none")
		}

		// Remove potentially dangerous headers
		c.Set("X-Powered-By", "")
		c.Set("Server", "")

		// Custom headers
		for name, value := range sh.config.CustomHeaders {
			c.Set(name, value)
		}

		// Security-related headers for APIs
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("X-Download-Options", "noopen")

		return c.Next()
	}
}

// DefaultHeadersConfig returns a secure default configuration
func DefaultHeadersConfig() HeadersConfig {
	return HeadersConfig{
		CSP:                  "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; manifest-src 'self';",
		HSTS:                 true,
		HSTSMaxAge:           31536000, // 1 year
		NoSniff:              true,
		FrameOptions:         "DENY",
		XSSProtection:        true,
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		PermittedCrossDomain: false,
		CustomHeaders:        make(map[string]string),
	}
}

// APIHeadersConfig returns headers configuration optimized for APIs
func APIHeadersConfig() HeadersConfig {
	return HeadersConfig{
		CSP:                  "default-src 'none'", // Very restrictive for APIs
		HSTS:                 true,
		HSTSMaxAge:           31536000,
		NoSniff:              true,
		FrameOptions:         "DENY",
		XSSProtection:        false, // Not relevant for APIs
		ReferrerPolicy:       "no-referrer",
		PermittedCrossDomain: false,
		CustomHeaders:        make(map[string]string),
	}
}

// WebAppHeadersConfig returns headers configuration optimized for web applications
func WebAppHeadersConfig() HeadersConfig {
	return HeadersConfig{
		CSP:                  "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'self'; frame-ancestors 'none'; form-action 'self'; base-uri 'self';",
		HSTS:                 true,
		HSTSMaxAge:           31536000,
		NoSniff:              true,
		FrameOptions:         "SAMEORIGIN", // Allow framing from same origin
		XSSProtection:        true,
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		PermittedCrossDomain: false,
		CustomHeaders:        make(map[string]string),
	}
}

// DevelopmentHeadersConfig returns headers configuration for development
func DevelopmentHeadersConfig() HeadersConfig {
	config := DefaultHeadersConfig()

	// More lenient CSP for development
	config.CSP = "default-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src 'self' data: https: http:; font-src 'self' https: http:; connect-src 'self' ws: wss: https: http:;"
	config.HSTS = false // Disable HSTS in development
	config.FrameOptions = "SAMEORIGIN"

	return config
}

// WithCustomHeader adds a custom header to the configuration
func (config *HeadersConfig) WithCustomHeader(name, value string) *HeadersConfig {
	if config.CustomHeaders == nil {
		config.CustomHeaders = make(map[string]string)
	}
	config.CustomHeaders[name] = value
	return config
}

// WithCSP sets the Content Security Policy
func (config *HeadersConfig) WithCSP(csp string) *HeadersConfig {
	config.CSP = csp
	return config
}

// WithHSTS configures HTTP Strict Transport Security
func (config *HeadersConfig) WithHSTS(enabled bool, maxAge int) *HeadersConfig {
	config.HSTS = enabled
	config.HSTSMaxAge = maxAge
	return config
}

// WithFrameOptions sets the X-Frame-Options header
func (config *HeadersConfig) WithFrameOptions(option string) *HeadersConfig {
	config.FrameOptions = option
	return config
}

// LogSecurityHeaders logs the applied security headers for debugging
func (sh *SecurityHeaders) LogSecurityHeaders(c *fiber.Ctx) {
	if sh.logger == nil {
		return
	}

	headers := []string{
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Referrer-Policy",
		"X-Permitted-Cross-Domain-Policies",
	}

	fields := []zap.Field{
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
	}

	for _, header := range headers {
		if value := c.Get(header); value != "" {
			fields = append(fields, zap.String(header, value))
		}
	}

	sh.logger.Debug("Applied security headers", fields...)
}
