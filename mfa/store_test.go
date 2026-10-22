package mfa

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_MarkCodeUsed(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	created, err := store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)

	created, err = store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.False(t, created)
}

func TestMemoryStore_AllowsReuseAfterExpiry(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	created, err := store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)

	store.mutex.Lock()
	store.entries["totp:used:user-1:1"] = time.Now().Add(-time.Second)
	store.mutex.Unlock()

	created, err = store.MarkCodeUsed(ctx, "totp:used:user-1:1", time.Minute)
	require.NoError(t, err)
	assert.True(t, created)
}
