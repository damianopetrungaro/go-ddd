package order

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("could not find order")
	ErrNotAdded = errors.New("could not add order")
)

// Repo represents the layer to read/write data from/to the storage
// You can also consider splitting this interface into multiple ones
type Repo interface {
	Get(ctx context.Context, id ID) (*Order, error)
	Add(ctx context.Context, order *Order) error
}
