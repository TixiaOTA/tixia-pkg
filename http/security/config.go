package security

import (
	"time"
)

// Config holds all security middleware configuration
type Config struct {
	RateLimit  *RateLimitConfig  `json:"rate_limit"`
	Headers    *HeadersConfig    `json:"headers"`
	CORS       *CORSConfig       `json:"cors"`
	IPFilter   *IPFilterConfig   `json:"ip_filter"`
	Validation *ValidationConfig `json:"validation"`
}

// RateLimitConfig configures rate limiting middleware
type RateLimitConfig struct {
	Enabled bool          `json:"enabled"`
	Rate    int           `json:"rate"`     // requests per window
	Window  time.Duration `json:"window"`   // time window duration
	Burst   int           `json:"burst"`    // burst capacity
	KeyFunc string        `json:"key_func"` // "ip", "user", "endpoint"
	Message string        `json:"message"`  // custom rate limit exceeded message
}

// HeadersConfig configures security headers middleware
type HeadersConfig struct {
	Enabled              bool              `json:"enabled"`
	CSP                  string            `json:"csp"`                    // Content Security Policy
	HSTS                 bool              `json:"hsts"`                   // HTTP Strict Transport Security
	HSTSMaxAge           int               `json:"hsts_max_age"`           // HSTS max age in seconds
	NoSniff              bool              `json:"no_sniff"`               // X-Content-Type-Options: nosniff
	FrameOptions         string            `json:"frame_options"`          // X-Frame-Options
	XSSProtection        bool              `json:"xss_protection"`         // X-XSS-Protection
	ReferrerPolicy       string            `json:"referrer_policy"`        // Referrer-Policy
	PermittedCrossDomain bool              `json:"permitted_cross_domain"` // X-Permitted-Cross-Domain-Policies
	CustomHeaders        map[string]string `json:"custom_headers"`         // Custom security headers
}

// CORSConfig configures CORS middleware with security considerations
type CORSConfig struct {
	Enabled             bool     `json:"enabled"`
	AllowOrigins        []string `json:"allow_origins"`         // Allowed origins
	AllowMethods        []string `json:"allow_methods"`         // Allowed HTTP methods
	AllowHeaders        []string `json:"allow_headers"`         // Allowed headers
	AllowCredentials    bool     `json:"allow_credentials"`     // Allow credentials
	ExposeHeaders       []string `json:"expose_headers"`        // Headers to expose
	MaxAge              int      `json:"max_age"`               // Max age for preflight requests
	AllowPrivateNetwork bool     `json:"allow_private_network"` // Allow private network access
	SecurityHeaders     bool     `json:"security_headers"`      // Add security headers to CORS responses
}

// IPFilterConfig configures IP filtering middleware
type IPFilterConfig struct {
	Enabled   bool     `json:"enabled"`
	Whitelist []string `json:"whitelist"` // CIDR blocks or IPs to allow
	Blacklist []string `json:"blacklist"` // CIDR blocks or IPs to deny
	Mode      string   `json:"mode"`      // "whitelist" or "blacklist"
	Message   string   `json:"message"`   // custom blocked message
}

// ValidationConfig configures request validation middleware
type ValidationConfig struct {
	Enabled             bool     `json:"enabled"`
	MaxBodySize         int64    `json:"max_body_size"`         // maximum request body size in bytes
	AllowedMethods      []string `json:"allowed_methods"`       // allowed HTTP methods
	RequiredHeaders     []string `json:"required_headers"`      // headers that must be present
	BlockedPaths        []string `json:"blocked_paths"`         // paths to block (regex patterns)
	AllowedContentTypes []string `json:"allowed_content_types"` // allowed content types
	ValidateJSON        bool     `json:"validate_json"`         // validate JSON payloads
	ValidateXML         bool     `json:"validate_xml"`          // validate XML payloads
}

// DefaultConfig returns a configuration with sensible defaults
// All middleware is disabled by default to avoid surprises
func DefaultConfig() *Config {
	return &Config{
		RateLimit: &RateLimitConfig{
			Enabled: false,
			Rate:    100,
			Window:  time.Minute,
			Burst:   10,
			KeyFunc: "ip",
			Message: "Rate limit exceeded. Please try again later.",
		},
		Headers: &HeadersConfig{
			Enabled:              false,
			CSP:                  "default-src 'self'",
			HSTS:                 true,
			HSTSMaxAge:           31536000, // 1 year
			NoSniff:              true,
			FrameOptions:         "DENY",
			XSSProtection:        true,
			ReferrerPolicy:       "strict-origin-when-cross-origin",
			PermittedCrossDomain: false,
			CustomHeaders:        make(map[string]string),
		},
		CORS: &CORSConfig{
			Enabled:             false,
			AllowOrigins:        []string{"http://localhost:3000", "https://localhost:3000"},
			AllowMethods:        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:        []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-Locale"},
			AllowCredentials:    true,
			ExposeHeaders:       []string{},
			MaxAge:              86400, // 24 hours
			AllowPrivateNetwork: false,
			SecurityHeaders:     true,
		},
		IPFilter: &IPFilterConfig{
			Enabled:   false,
			Whitelist: []string{},
			Blacklist: []string{},
			Mode:      "blacklist",
			Message:   "Access denied from your IP address.",
		},
		Validation: &ValidationConfig{
			Enabled:             false,
			MaxBodySize:         10 << 20, // 10MB
			AllowedMethods:      []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			RequiredHeaders:     []string{},
			BlockedPaths:        []string{},
			AllowedContentTypes: []string{"application/json", "application/xml", "text/plain", "multipart/form-data"},
			ValidateJSON:        true,
			ValidateXML:         false,
		},
	}
}

// ProductionDefaults returns production-ready defaults with security features enabled
func ProductionDefaults() *Config {
	config := DefaultConfig()

	// Enable headers with production settings
	config.Headers.Enabled = true

	// Enable basic validation
	config.Validation.Enabled = true
	config.Validation.MaxBodySize = 5 << 20 // 5MB for production

	return config
}

// DevelopmentDefaults returns development-friendly defaults
func DevelopmentDefaults() *Config {
	config := DefaultConfig()

	// More lenient settings for development
	config.Validation.Enabled = true
	config.Validation.MaxBodySize = 50 << 20 // 50MB for development

	return config
}
