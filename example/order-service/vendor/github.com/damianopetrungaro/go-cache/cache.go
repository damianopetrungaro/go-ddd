package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// List of errors returned by the Cache implementations
var (
	ErrNotSet    = errors.New("could not set cache value")
	ErrNotGet    = errors.New("could not get cache value")
	ErrNotFound  = fmt.Errorf("%w: could not find cache value", ErrNotGet)
	ErrExpired   = fmt.Errorf("%w: could not get expired cache value", ErrNotGet)
	ErrNotDelete = errors.New("could not delete cache value")
)

const (
	// NoExpiration is a constant used to mark an item as never expiring
	NoExpiration = time.Duration(0)
)

// Cache represents the contract for interacting with a cache layer
type Cache[K comparable, V any] interface {
	Get(context.Context, K) (V, error)
	Set(context.Context, K, V, time.Duration) error
	Delete(context.Context, K) error
}
