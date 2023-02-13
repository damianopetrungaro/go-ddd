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
type Repo interface {
	Find(ctx context.Context, id ID) (*Order, error)
	Add(ctx context.Context, order *Order) error
}
