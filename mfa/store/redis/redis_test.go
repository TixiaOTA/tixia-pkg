package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/TixiaOTA/tixia-pkg/redis/cache"
	redisclient "github.com/TixiaOTA/tixia-pkg/redis/client"
)

func newTestStore(t *testing.T) (*Store, *miniredis.Miniredis) {
	t.Helper()

	server := miniredis.RunT(t)

	cacheManager, err := cache.NewCacheManager(cache.CacheConfig{
		BaseRedisConfig: redisclient.BaseRedisConfig{
			URL: "redis://" + server.Addr(),
		},
		DefaultTTL: time.Minute,
		MaxTTL:     time.Hour,
		KeyPrefix:  "mfa-test",
	}, logger.New(logger.Config{Level: "error", Environment: "development"}))
	require.NoError(t, err)

	t.Cleanup(func() {
		cacheManager.Close()
	})

	store, err := New(cacheManager)
	require.NoError(t, err)

	return store, server
}

func TestNew_NilCacheManager(t *testing.T) {
	_, err := New(nil)
	assert.ErrorIs(t, err, ErrNilCacheManager)
}

func TestStore_MarkCodeUsed(t *testing.T) {
	store, _ := newTestStore(t)
	ctx := context.Background()

	created, err := store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)

	created, err = store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.False(t, created)
}

func TestStore_MarkCodeUsedAfterExpiry(t *testing.T) {
	store, server := newTestStore(t)
	ctx := context.Background()

	created, err := store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)

	server.FastForward(2 * time.Minute)

	created, err = store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)
}
