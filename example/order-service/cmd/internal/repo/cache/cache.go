package cache

import (
	"context"
	"github.com/damianopetrungaro/go-cache"
	"github.com/damianopetrungaro/golog"
	"github.com/organization/order-service"
	"time"
)

var (
	_ order.Repo = &Cache{}

	// memory represents the default in-memory cache
	memory = cache.NewInMemory[order.ID, *order.Order](time.Minute, 10_000)

	// defaultTTL represent the TTL for an item in the cache
	defaultTTL = time.Minute
)

// Cache represents a cache layer for the order.Repo
type Cache struct {
	base   order.Repo
	store  cache.Cache[order.ID, *order.Order]
	logger golog.Logger
}

// DefaultStore returns a default cache store
func DefaultStore() cache.Cache[order.ID, *order.Order] {
	return memory
}

// New returns a cache wrapper for an order.Repo
func New(base order.Repo, store cache.Cache[order.ID, *order.Order], logger golog.Logger) *Cache {
	return &Cache{base: base, store: store, logger: logger}
}

// Get tries getting an order from the cache storage first
// if not found calls the Find method of the base
func (c *Cache) Get(ctx context.Context, id order.ID) (*order.Order, error) {
	o, err := c.store.Get(ctx, id)
	switch err {
	case nil:
		c.logger.With(golog.Err(err)).Debug(ctx, "order was found in cache")
		return o, nil
	default:
		c.logger.With(golog.Err(err)).Debug(ctx, "order was not found in cache")
		o, err := c.base.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		c.sets(ctx, o)
		return o, nil
	}
}

// Add sets an order in the cache if successfully added
func (c *Cache) Add(ctx context.Context, o *order.Order) error {
	if err := c.base.Add(ctx, o); err != nil {
		c.logger.With(golog.Err(err)).Warn(ctx, "order was not added in cache")
		return err
	}

	c.sets(ctx, o)

	return nil
}

func (c *Cache) sets(ctx context.Context, o *order.Order) {
	if err := c.store.Set(ctx, o.ID, o, defaultTTL); err != nil {
		c.logger.With(golog.Err(err)).Warn(ctx, "order was not set in cache")
	}

	c.logger.Info(ctx, "order was set in cache")
}
