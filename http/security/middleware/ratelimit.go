package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiter implements rate limiting using Redis and token bucket algorithm
type RateLimiter struct {
	redis  *redis.Client
	logger *zap.Logger
	config RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Rate    int           // requests per window
	Window  time.Duration // time window duration
	Burst   int           // burst capacity
	KeyFunc string        // "ip", "user", "endpoint"
	Message string        // custom rate limit exceeded message
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(redis *redis.Client, logger *zap.Logger, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis:  redis,
		logger: logger,
		config: config,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if Redis is not available
		if rl.redis == nil {
			if rl.logger != nil {
				rl.logger.Warn("Rate limiting skipped: Redis client not available")
			}
			return c.Next()
		}

		// Generate rate limiting key
		key := rl.generateKey(c)

		// Check rate limit
		allowed, remaining, resetTime, err := rl.checkRateLimit(c.Context(), key)
		if err != nil {
			if rl.logger != nil {
				rl.logger.Error("Rate limit check failed",
					zap.String("key", key),
					zap.Error(err),
				)
			}
			// On error, allow the request but log the issue
			return c.Next()
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.config.Rate))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		// IP
		ip := c.Get("X-Real-IP")
		if ip == "" {
			ip = c.IP()
		}

		// Check if rate limit exceeded
		if !allowed {
			if rl.logger != nil {
				rl.logger.Warn("Rate limit exceeded",
					zap.String("key", key),
					zap.String("ip", ip),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
			}

			c.Set("Retry-After", strconv.FormatInt(resetTime-time.Now().UTC().Unix(), 10))

			message := rl.config.Message
			if message == "" {
				message = "Rate limit exceeded. Please try again later."
			}

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "rate_limit_exceeded",
				"message":     message,
				"retry_after": resetTime - time.Now().UTC().Unix(),
			})
		}

		return c.Next()
	}
}

// generateKey creates a rate limiting key based on the configuration
func (rl *RateLimiter) generateKey(c *fiber.Ctx) string {
	ip := c.Get("X-Real-IP")
	if ip == "" {
		ip = c.IP()

	}
	switch rl.config.KeyFunc {
	case "user":
		userID := c.Locals("user_id")
		if userID != nil {
			return fmt.Sprintf("ratelimit:user:%v", userID)
		}
		// Fallback to IP if user ID not available
		return fmt.Sprintf("ratelimit:ip:%s", ip)
	case "endpoint":
		return fmt.Sprintf("ratelimit:endpoint:%s:%s", c.Method(), c.Path())
	case "user_endpoint":
		userID := c.Locals("user_id")
		if userID != nil {
			return fmt.Sprintf("ratelimit:user_endpoint:%v:%s:%s", userID, c.Method(), c.Path())
		}
		// Fallback to IP+endpoint if user ID not available
		return fmt.Sprintf("ratelimit:ip_endpoint:%s:%s:%s", ip, c.Method(), c.Path())
	default: // "ip"
		return fmt.Sprintf("ratelimit:ip:%s", ip)
	}
}

// checkRateLimit implements sliding window rate limiting using Redis
func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string) (allowed bool, remaining int, resetTime int64, err error) {
	now := time.Now().UTC()
	window := rl.config.Window
	windowStart := now.Truncate(window)

	// Redis Lua script for atomic rate limiting check
	script := `
		local key = KEYS[1]
		local window_start = ARGV[1]
		local window_size = ARGV[2]
		local limit = tonumber(ARGV[3])
		local burst = tonumber(ARGV[4])
		local now = tonumber(ARGV[5])
		
		-- Get current count
		local current = redis.call('GET', key)
		if current == false then
			current = 0
		else
			current = tonumber(current)
		end
		
		-- Check if we're within limits
		if current < limit then
			-- Increment counter
			local new_count = redis.call('INCR', key)
			
			-- Set expiration on first increment
			if new_count == 1 then
				redis.call('EXPIRE', key, window_size)
			end
			
			return {1, limit - new_count, window_start + window_size}
		else
			-- Rate limit exceeded
			return {0, 0, window_start + window_size}
		end
	`

	result, err := rl.redis.Eval(ctx, script, []string{key},
		windowStart.Unix(),
		int(window.Seconds()),
		rl.config.Rate,
		rl.config.Burst,
		now.Unix(),
	).Result()

	if err != nil {
		return false, 0, 0, fmt.Errorf("rate limit check failed: %w", err)
	}

	// Parse result
	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) != 3 {
		return false, 0, 0, fmt.Errorf("unexpected Redis response format")
	}

	allowedInt, _ := resultSlice[0].(int64)
	remainingInt, _ := resultSlice[1].(int64)
	resetTimeInt, _ := resultSlice[2].(int64)

	return allowedInt == 1, int(remainingInt), resetTimeInt, nil
}

// Reset clears the rate limit for a specific key (useful for testing or admin override)
func (rl *RateLimiter) Reset(ctx context.Context, key string) error {
	if rl.redis == nil {
		return fmt.Errorf("Redis client not available")
	}

	return rl.redis.Del(ctx, key).Err()
}

// GetStats returns current rate limit statistics for a key
func (rl *RateLimiter) GetStats(ctx context.Context, key string) (current int, ttl time.Duration, err error) {
	if rl.redis == nil {
		return 0, 0, fmt.Errorf("Redis client not available")
	}

	// Get current count
	currentStr, err := rl.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, 0, nil // Key doesn't exist
	}
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get current count: %w", err)
	}

	current, err = strconv.Atoi(currentStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid current count: %w", err)
	}

	// Get TTL
	ttl, err = rl.redis.TTL(ctx, key).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return current, ttl, nil
}
