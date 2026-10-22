package redis

import (
	"context"
	"errors"
	"time"

	"github.com/TixiaOTA/tixia-pkg/redis/cache"
)

// ErrNilCacheManager is returned when no cache manager is provided.
var ErrNilCacheManager = errors.New("mfa/redis: cache manager is required")

// Store is a Redis-backed mfa.Store implementation.
type Store struct {
	cache *cache.CacheManager
}

// New creates a new Redis-backed replay store.
func New(cacheManager *cache.CacheManager) (*Store, error) {
	if cacheManager == nil {
		return nil, ErrNilCacheManager
	}

	return &Store{
		cache: cacheManager,
	}, nil
}

// MarkCodeUsed atomically records a code as used until its TTL expires.
func (s *Store) MarkCodeUsed(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return s.cache.SetNX(ctx, key, true, &cache.CacheOptions{
		TTL: ttl,
	})
}
