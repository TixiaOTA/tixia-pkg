package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/TixiaOTA/tixia-pkg/redis/client"
)

func newTestCacheManager(t *testing.T) (*CacheManager, *miniredis.Miniredis) {
	t.Helper()

	server := miniredis.RunT(t)

	manager, err := NewCacheManager(CacheConfig{
		BaseRedisConfig: client.BaseRedisConfig{
			URL: "redis://" + server.Addr(),
		},
		DefaultTTL: time.Minute,
		MaxTTL:     time.Hour,
		KeyPrefix:  "test",
	}, logger.New(logger.Config{Level: "error", Environment: "development"}))
	require.NoError(t, err)

	t.Cleanup(func() {
		manager.Close()
	})

	return manager, server
}

func TestCacheManager_SetNX(t *testing.T) {
	manager, _ := newTestCacheManager(t)
	ctx := context.Background()

	created, err := manager.SetNX(ctx, "one-time", true, nil)
	require.NoError(t, err)
	assert.True(t, created)

	created, err = manager.SetNX(ctx, "one-time", true, nil)
	require.NoError(t, err)
	assert.False(t, created)

	var value bool
	found, err := manager.Get(ctx, "one-time", &value, nil)
	require.NoError(t, err)
	assert.True(t, found)
	assert.True(t, value)
}

func TestCacheManager_SetNX_Expiry(t *testing.T) {
	manager, server := newTestCacheManager(t)
	ctx := context.Background()

	created, err := manager.SetNX(ctx, "one-time", true, &CacheOptions{
		TTL: time.Minute,
	})
	require.NoError(t, err)
	assert.True(t, created)

	server.FastForward(2 * time.Minute)

	created, err = manager.SetNX(ctx, "one-time", true, &CacheOptions{
		TTL: time.Minute,
	})
	require.NoError(t, err)
	assert.True(t, created)
}
