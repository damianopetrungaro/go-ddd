package instrument

import (
	"context"
	"github.com/organization/order-service"
	"github.com/organization/order-service/internal"
)

var (
	_ Service = &internal.Service{}
)

// Service represents an interface matching internal.Service
// it is used as base dor generating code for instrument purposes
type Service interface {
	Place(context.Context, order.Number, order.UserID) (*order.Order, error)
	MarkAsShipped(context.Context, order.ID) (*order.Order, error)
	MarkAsDelivered(context.Context, order.ID) (*order.Order, error)
}

// NewService returns an instrumented Service
func NewService(base Service, name string) Service {
	return NewServiceWithPrometheus(NewServiceWithTracing(base, name), name)
}
