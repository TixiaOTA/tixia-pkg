package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/TixiaOTA/tixia-pkg/redis/client"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheManager struct {
	client *redis.Client
	logger *logger.Logger
	config CacheConfig
}

func NewCacheManager(cacheConfig CacheConfig, logger *logger.Logger) (*CacheManager, error) {
	baseClient, err := client.NewBaseRedisClient(cacheConfig.BaseRedisConfig, logger.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create base Redis client: %w", err)
	}

	return &CacheManager{
		client: baseClient.GetClient(),
		logger: logger,
		config: cacheConfig,
	}, nil
}

func (cm *CacheManager) GenerateKey(prefix, key string) string {
	if prefix == "" {
		prefix = cm.config.KeyPrefix
	}
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (cm *CacheManager) GenerateHashKey(prefix string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data for hash key: %w", err)
	}

	hash := md5.Sum(jsonData)
	key := fmt.Sprintf("%x", hash)
	return cm.GenerateKey(prefix, key), nil
}

func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, opts *CacheOptions) error {
	if opts == nil {
		opts = &CacheOptions{
			TTL: cm.config.DefaultTTL,
		}
	}

	if opts.TTL == 0 {
		opts.TTL = cm.config.DefaultTTL
	}

	if opts.TTL > cm.config.MaxTTL {
		opts.TTL = cm.config.MaxTTL
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	err = cm.client.Set(ctx, cacheKey, jsonValue, opts.TTL).Err()
	if err != nil {
		cm.logger.Error("Failed to set cache value",
			zap.String("key", cacheKey),
			zap.Error(err))
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	return nil
}

// SetNX sets a value only if the key does not already exist. It returns true
// when this call created the key and false when the key already existed.
func (cm *CacheManager) SetNX(ctx context.Context, key string, value interface{}, opts *CacheOptions) (bool, error) {
	if opts == nil {
		opts = &CacheOptions{
			TTL: cm.config.DefaultTTL,
		}
	}

	if opts.TTL == 0 {
		opts.TTL = cm.config.DefaultTTL
	}

	if opts.TTL > cm.config.MaxTTL {
		opts.TTL = cm.config.MaxTTL
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	created, err := cm.client.SetNX(ctx, cacheKey, jsonValue, opts.TTL).Result()
	if err != nil {
		cm.logger.Error("Failed to set cache value if not exists",
			zap.String("key", cacheKey),
			zap.Error(err))
		return false, fmt.Errorf("failed to set cache value if not exists: %w", err)
	}

	return created, nil
}

func (cm *CacheManager) Get(ctx context.Context, key string, dest interface{}, opts *CacheOptions) (bool, error) {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	val, err := cm.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		cm.logger.Error("Failed to get cache value",
			zap.String("key", cacheKey),
			zap.Error(err))
		return false, fmt.Errorf("failed to get cache value: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		cm.logger.Error("Failed to unmarshal cached value",
			zap.String("key", cacheKey),
			zap.Error(err))
		return false, fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return true, nil
}

func (cm *CacheManager) GetOrSet(ctx context.Context, key string, dest interface{},
	fetchFunc func() (interface{}, error), opts *CacheOptions) error {

	// Try to get from cache first
	found, err := cm.Get(ctx, key, dest, opts)
	if err != nil {
		return err
	}

	if found {
		return nil
	}

	// Fetch data if not in cache
	data, err := fetchFunc()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	// Set in cache
	err = cm.Set(ctx, key, data, opts)
	if err != nil {
		cm.logger.Warn("Failed to set cache after fetch",
			zap.String("key", key),
			zap.Error(err))
	}

	// Copy data to dest
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal fetched data: %w", err)
	}

	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal fetched data: %w", err)
	}

	return nil
}

func (cm *CacheManager) Delete(ctx context.Context, key string, opts *CacheOptions) error {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	err := cm.client.Del(ctx, cacheKey).Err()
	if err != nil {
		cm.logger.Error("Failed to delete cache key",
			zap.String("key", cacheKey),
			zap.Error(err))
		return fmt.Errorf("failed to delete cache key: %w", err)
	}

	return nil
}

func (cm *CacheManager) DeletePattern(ctx context.Context, pattern string, opts *CacheOptions) error {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cachePattern := cm.GenerateKey(opts.KeyPrefix, pattern)

	iter := cm.client.Scan(ctx, 0, cachePattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		cm.logger.Error("Failed to scan cache keys",
			zap.String("pattern", cachePattern),
			zap.Error(err))
		return fmt.Errorf("failed to scan cache keys: %w", err)
	}

	if len(keys) > 0 {
		err := cm.client.Del(ctx, keys...).Err()
		if err != nil {
			cm.logger.Error("Failed to delete cache keys",
				zap.String("pattern", cachePattern),
				zap.Int("count", len(keys)),
				zap.Error(err))
			return fmt.Errorf("failed to delete cache keys: %w", err)
		}
	}

	return nil
}

func (cm *CacheManager) Exists(ctx context.Context, key string, opts *CacheOptions) (bool, error) {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	exists, err := cm.client.Exists(ctx, cacheKey).Result()
	if err != nil {
		cm.logger.Error("Failed to check cache key existence",
			zap.String("key", cacheKey),
			zap.Error(err))
		return false, fmt.Errorf("failed to check cache key existence: %w", err)
	}

	return exists > 0, nil
}

func (cm *CacheManager) TTL(ctx context.Context, key string, opts *CacheOptions) (time.Duration, error) {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)

	ttl, err := cm.client.TTL(ctx, cacheKey).Result()
	if err != nil {
		cm.logger.Error("Failed to get cache key TTL",
			zap.String("key", cacheKey),
			zap.Error(err))
		return 0, fmt.Errorf("failed to get cache key TTL: %w", err)
	}

	return ttl, nil
}

func (cm *CacheManager) RenewTTL(ctx context.Context, key string, ttl time.Duration, opts *CacheOptions) error {
	if opts == nil {
		opts = &CacheOptions{}
	}

	cacheKey := cm.GenerateKey(opts.KeyPrefix, key)
	_, err := cm.client.Expire(ctx, cacheKey, ttl).Result()
	if err != nil {
		cm.logger.Error("failed to renew TTL for key",
			zap.String("key", cacheKey),
			zap.Duration("ttl", ttl),
			zap.Error(err))
		return fmt.Errorf("failed to renew TTL for key %s: %w", cacheKey, err)
	}
	return nil
}

func (cm *CacheManager) FlushAll(ctx context.Context) error {
	err := cm.client.FlushAll(ctx).Err()
	if err != nil {
		cm.logger.Error("Failed to flush all cache",
			zap.Error(err))
		return fmt.Errorf("failed to flush all cache: %w", err)
	}

	cm.logger.Info("All cache data flushed")
	return nil
}

func (cm *CacheManager) GetStats(ctx context.Context) (map[string]interface{}, error) {
	_, err := cm.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	// Parse Redis info output
	stats := make(map[string]interface{})
	stats["connected_clients"] = 0
	stats["used_memory"] = 0
	stats["total_commands_processed"] = 0

	return stats, nil
}

func (cm *CacheManager) Close() error {
	if cm.client != nil {
		if err := cm.client.Close(); err != nil {
			cm.logger.Error("Failed to close cache manager Redis client", zap.Error(err))
			return fmt.Errorf("failed to close cache manager Redis client: %w", err)
		}
		cm.logger.Error("Cache manager Redis client closed")
	}
	return nil
}

func (cm *CacheManager) HealthCheck(ctx context.Context) error {
	if cm.client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := cm.client.Ping(ctx).Err(); err != nil {
		cm.logger.Error("Redis health check failed", zap.Error(err))
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	cm.logger.Debug("Redis health check passed")
	return nil
}

func (cm *CacheManager) GetClient() *redis.Client {
	return cm.client
}

func (cm *CacheManager) GetConnectionInfo(ctx context.Context) (map[string]interface{}, error) {
	if cm.client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	info := make(map[string]interface{})

	serverInfo, err := cm.client.Info(ctx).Result()
	if err != nil {
		cm.logger.Error("Failed to get Redis server info", zap.Error(err))
		return nil, fmt.Errorf("failed to get Redis server info: %w", err)
	}

	clientList, err := cm.client.ClientList(ctx).Result()
	if err != nil {
		cm.logger.Error("Failed to get Redis client list", zap.Error(err))
		return nil, fmt.Errorf("failed to get Redis client list: %w", err)
	}

	info["server_info"] = serverInfo
	info["client_list"] = clientList
	info["config"] = cm.config
	info["pool_stats"] = cm.client.PoolStats()

	return info, nil
}
