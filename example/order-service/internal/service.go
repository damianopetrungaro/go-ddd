package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/damianopetrungaro/golog"
	"github.com/organization/order-service"
)

// Errors that the application layer exposes
var (
	ErrNotPlaced            = errors.New("order could not be placed")
	ErrNotMarkedAsShipped   = errors.New("order could not be marked as shipped")
	ErrNotMarkedAsDelivered = errors.New("order could not be marked as delivered")
)

// Service represent the application layer
// it depends on the domain logic and can be used by any infrastructure layer as domain logic orchestrator
type Service struct {
	repo   order.Repo
	logger golog.Logger
}

// NewService returns a new Service
func NewService(repo order.Repo, logger golog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Place places an order and store it in the repository
func (s *Service) Place(ctx context.Context, n order.Number, uID order.UserID) (*order.Order, error) {
	o := order.Place(n, uID)

	if err := s.repo.Add(ctx, o); err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not added once placed")
		return nil, fmt.Errorf("%w: %w", ErrNotPlaced, err)
	}

	return o, nil
}

// MarkAsShipped marks as shipped an order and store it in the repository
func (s *Service) MarkAsShipped(ctx context.Context, id order.ID) (*order.Order, error) {
	o, err := s.repo.Find(ctx, id)
	if err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not found")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsShipped, err)
	}

	if err := o.MarkAsShipped(); err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not marked as shipped")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsShipped, err)
	}

	if err := s.repo.Add(ctx, o); err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not added once marked as shipped")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsShipped, err)
	}

	return o, nil
}

// MarkAsDelivered marks as delivered an order and store it in the repository
func (s *Service) MarkAsDelivered(ctx context.Context, id order.ID) (*order.Order, error) {
	o, err := s.repo.Find(ctx, id)
	if err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not found")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsDelivered, err)
	}

	if err := o.MarkAsDelivered(); err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not marked as delivered")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsDelivered, err)
	}

	if err := s.repo.Add(ctx, o); err != nil {
		s.logger.With(golog.Err(err)).Error(ctx, "order was not added once marked as delivered")
		return nil, fmt.Errorf("%w: %w", ErrNotMarkedAsDelivered, err)
	}

	return o, nil
}
