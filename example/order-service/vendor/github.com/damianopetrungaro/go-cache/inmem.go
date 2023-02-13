package cache

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

var _ Cache[string, any] = &InMem[string, any]{}

type expiresAt int64

func (ea expiresAt) isExpired() bool {
	i := int64(ea)
	return time.Now().UnixNano() > i && i != int64(NoExpiration)
}

type item[V any] struct {
	val       V
	expiresAt expiresAt
}

// InMem is a Cache implementation which interacts with an in-memory map
// It is concurrent safe
type InMem[K comparable, V any] struct {
	items  map[K]item[V]
	cap    int
	ticker *time.Ticker
	mu     sync.RWMutex
}

// NewInMemory returns a InMem instance
func NewInMemory[K comparable, V any](cleanUpInterval time.Duration, cap int) *InMem[K, V] {
	inmem := &InMem[K, V]{
		items:  map[K]item[V]{},
		cap:    cap,
		ticker: time.NewTicker(cleanUpInterval),
	}

	go func() {
		for range inmem.ticker.C {
			inmem.mu.Lock()
			inmem.cleanup()
			inmem.mu.Unlock()
		}
	}()

	return inmem
}

// Get retrieves an item from an in-memory map
func (i *InMem[K, V]) Get(ctx context.Context, key K) (V, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	select {
	case <-ctx.Done():
		return *new(V), fmt.Errorf("%w: %s", ErrNotGet, ctx.Err())
	default:
	}

	item, ok := i.items[key]
	if !ok {
		return *new(V), ErrNotFound
	}

	if item.expiresAt.isExpired() {
		return *new(V), ErrExpired
	}

	return item.val, nil
}

// Set stores an item to an in-memory map
func (i *InMem[K, V]) Set(ctx context.Context, key K, val V, ttl time.Duration) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", ErrNotSet, ctx.Err())
	default:
	}

	exp := expiresAt(ttl)
	if ttl != NoExpiration {
		exp = expiresAt(time.Now().Add(ttl).UnixNano())
	}

	if len(i.items) == i.cap {
		i.cleanup()
	}

	i.items[key] = item[V]{val: val, expiresAt: exp}

	return nil
}

// Delete removes an item to an in-memory map
func (i *InMem[K, V]) Delete(ctx context.Context, key K) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", ErrNotDelete, ctx.Err())
	default:
	}

	delete(i.items, key)
	return nil
}

// Close stops the inner ticker
func (i *InMem[K, V]) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.ticker.Stop()
	return nil
}

// cleanup remove all the expired items.
// if no item is expired, it deletes the one closer to expire
func (i *InMem[K, V]) cleanup() {
	ks := []K{}
	minExp := math.MaxInt64

	for k, item := range i.items {
		switch {
		case item.expiresAt.isExpired():
			delete(i.items, k)
		case minExp == int(item.expiresAt):
			minExp = int(item.expiresAt)
			ks = append(ks, k)
		case minExp > int(item.expiresAt):
			minExp = int(item.expiresAt)
			ks = []K{k}
		}
	}

	if len(i.items) < i.cap {
		return
	}

	for _, k := range ks {
		delete(i.items, k)
	}
}
