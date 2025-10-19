package client

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type BaseRedisClient struct {
	client *redis.Client
	logger *zap.Logger
	config BaseRedisConfig
}

func NewBaseRedisClient(config BaseRedisConfig, logger *zap.Logger) (*BaseRedisClient, error) {
	opts, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}

	// Apply configuration overrides
	if config.PoolSize > 0 {
		opts.PoolSize = config.PoolSize
	}
	if config.MinIdleConns > 0 {
		opts.MinIdleConns = config.MinIdleConns
	}
	if config.MaxRetries > 0 {
		opts.MaxRetries = config.MaxRetries
	}
	if config.DialTimeout > 0 {
		opts.DialTimeout = config.DialTimeout
	}
	if config.ReadTimeout > 0 {
		opts.ReadTimeout = config.ReadTimeout
	}
	if config.WriteTimeout > 0 {
		opts.WriteTimeout = config.WriteTimeout
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &BaseRedisClient{
		client: client,
		logger: logger,
		config: config,
	}, nil
}

func (b *BaseRedisClient) GetClient() *redis.Client {
	return b.client
}

func (b *BaseRedisClient) HealthCheck(ctx context.Context) error {
	if b.client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := b.client.Ping(ctx).Err(); err != nil {
		b.logger.Error("Redis health check failed", zap.Error(err))
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	b.logger.Debug("Redis health check passed")
	return nil
}

func (b *BaseRedisClient) Close() error {
	if b.client != nil {
		if err := b.client.Close(); err != nil {
			b.logger.Error("Failed to close Redis connection", zap.Error(err))
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
		b.logger.Error("Redis connection closed")
	}
	return nil
}

func (b *BaseRedisClient) GetConnectionInfo(ctx context.Context) (map[string]interface{}, error) {
	if b.client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	info := make(map[string]interface{})

	serverInfo, err := b.client.Info(ctx).Result()
	if err != nil {
		b.logger.Error("Failed to get Redis server info", zap.Error(err))
		return nil, fmt.Errorf("failed to get Redis server info: %w", err)
	}

	clientList, err := b.client.ClientList(ctx).Result()
	if err != nil {
		b.logger.Error("Failed to get Redis client list", zap.Error(err))
		return nil, fmt.Errorf("failed to get Redis client list: %w", err)
	}

	info["server_info"] = serverInfo
	info["client_list"] = clientList
	info["config"] = b.config
	info["pool_stats"] = b.client.PoolStats()

	return info, nil
}

func (b *BaseRedisClient) GetStats(ctx context.Context) (map[string]interface{}, error) {
	if b.client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	stats := make(map[string]interface{})

	poolStats := b.client.PoolStats()
	stats["pool_stats"] = map[string]interface{}{
		"hits":     poolStats.Hits,
		"misses":   poolStats.Misses,
		"timeouts": poolStats.Timeouts,
	}

	info, err := b.client.Info(ctx, "stats").Result()
	if err != nil {
		b.logger.Error("Failed to get Redis stats", zap.Error(err))
		return nil, fmt.Errorf("failed to get Redis stats: %w", err)
	}

	stats["redis_info"] = info
	stats["config"] = b.config

	return stats, nil
}
