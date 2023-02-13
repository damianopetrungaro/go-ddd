package cache

import (
	"context"
	"time"
)

var (
	// DefaultMultiLevelExpiration is a constant used to mark an item to expire to the default value passed to a MultiLevel
	DefaultMultiLevelExpiration = time.Duration(-1)
)

// MultiLevel is a Cache implementation which allow a multi level usage cache
type MultiLevel[K comparable, V any] struct {
	local            *InMem[K, V]
	defaultLocalTTL  time.Duration
	remote           Cache[K, V]
	defaultRemoteTTL time.Duration
}

// NewMultiLevel returns a MultiLevel
func NewMultiLevel[K comparable, V any](
	local *InMem[K, V],
	defaultLocalTTL time.Duration,
	remote Cache[K, V],
	defaultRemoteTTL time.Duration,
) *MultiLevel[K, V] {
	return &MultiLevel[K, V]{
		local:            local,
		defaultLocalTTL:  defaultLocalTTL,
		remote:           remote,
		defaultRemoteTTL: defaultRemoteTTL,
	}
}

// Get search in local cache first, if an error occurred moves to the remote one
func (m *MultiLevel[K, V]) Get(ctx context.Context, k K) (V, error) {
	val, err := m.local.Get(ctx, k)
	if err != nil {
		return m.remote.Get(ctx, k)
	}
	return val, nil
}

// Set traverse all the caches, if all of them fail it returns a generic ErrNotSet
func (m *MultiLevel[K, V]) Set(ctx context.Context, k K, v V, ttl time.Duration) error {
	if ttl == DefaultMultiLevelExpiration {
		ttl = m.defaultRemoteTTL
	}
	if err := m.remote.Set(ctx, k, v, ttl); err != nil {
		return err
	}

	if ttl == DefaultMultiLevelExpiration {
		ttl = m.defaultLocalTTL
	}
	// in memory won't fail apart from having a ctx done.
	// passing a background to prevent it
	_ = m.local.Set(context.Background(), k, v, ttl)
	return nil
}

// Delete traverse all the caches, if all of them fail it returns a generic ErrNotDelete
func (m *MultiLevel[K, V]) Delete(ctx context.Context, k K) error {
	if err := m.remote.Delete(ctx, k); err != nil {
		return err
	}

	// in memory won't fail apart from having a ctx done.
	// passing a background to prevent it
	_ = m.local.Delete(context.Background(), k)
	return nil
}
