package mfa

import (
	"context"
	"sync"
	"time"
)

// Store tracks one-time codes so they cannot be replayed.
type Store interface {
	// MarkCodeUsed atomically records a code as used. It returns true when this
	// call recorded the code and false when it had already been recorded.
	MarkCodeUsed(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

// MemoryStore is an in-process Store implementation. It is suitable for tests
// and single-instance deployments. Use a shared store such as Redis when more
// than one instance verifies codes.
type MemoryStore struct {
	mutex   sync.Mutex
	entries map[string]time.Time
}

// NewMemoryStore creates a new in-memory replay store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make(map[string]time.Time),
	}
}

// MarkCodeUsed records a code as used until its TTL expires.
func (s *MemoryStore) MarkCodeUsed(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()

	if expiresAt, exists := s.entries[key]; exists && expiresAt.After(now) {
		return false, nil
	}

	for entryKey, expiresAt := range s.entries {
		if !expiresAt.After(now) {
			delete(s.entries, entryKey)
		}
	}

	s.entries[key] = now.Add(ttl)

	return true, nil
}
