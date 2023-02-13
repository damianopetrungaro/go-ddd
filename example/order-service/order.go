package order

import (
	"errors"
	"fmt"
	"time"
)

// Domain errors raised by the order aggregate
var (
	ErrNotShipped   = errors.New("could not mark the order as shipped")
	ErrNotDelivered = errors.New("could not mark the order as delivered")
)

// Order represents an order aggregate
type Order struct {
	ID          ID
	Number      Number
	Status      Status
	PlacedBy    UserID
	PlacedAt    time.Time
	ShippedAt   time.Time
	DeliveredAt time.Time
}

// Place places a new order
// It is a factory function that uses the ubiquitous language of the domain
func Place(number Number, placedBy UserID) *Order {
	return &Order{
		ID:       NewID(),
		Number:   number,
		Status:   Placed,
		PlacedBy: placedBy,
		PlacedAt: time.Now(),
	}
}

// MarkAsShipped marks an order as shipped
// It returns ErrNotShipped when the operation violated the domain invariants
func (o *Order) MarkAsShipped() error {
	switch {
	case o.PlacedAt.IsZero():
		return fmt.Errorf("%w: not placed", ErrNotShipped)
	case o.Status == Shipped:
		return fmt.Errorf("%w: already shipped", ErrNotShipped)
	case o.Status == Delivered:
		return fmt.Errorf("%w: already delivered", ErrNotShipped)
	}

	o.Status = Shipped
	o.ShippedAt = time.Now()
	return nil
}

// MarkAsDelivered marks an order as delivered
// It returns ErrNotDelivered when the operation violated the domain invariants
func (o *Order) MarkAsDelivered() error {
	switch {
	case o.ShippedAt.IsZero():
		return fmt.Errorf("%w: not shipped", ErrNotDelivered)
	case o.Status == Delivered:
		return fmt.Errorf("%w: already delivered", ErrNotDelivered)
	}

	o.Status = Delivered
	o.DeliveredAt = time.Now()
	return nil
}
