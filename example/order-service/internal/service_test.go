package internal

import (
	"context"
	"errors"
	gologTest "github.com/damianopetrungaro/golog/test"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/organization/order-service"
	"testing"
	"time"
)

func TestService_Place(t *testing.T) {
	t.Run("not added", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Add(ctx, gomock.Any()).Return(order.ErrNotAdded)

		svc := NewService(repo, logger)
		n := newOrderNumber(t)
		uID := newUserID(t)

		o, err := svc.Place(ctx, n, uID)
		if !errors.Is(err, ErrNotPlaced) {
			t.Fatalf("could not match placing order error: %s", err)
		}

		if o != nil {
			t.Fatalf("could not match a nil order: %v", o)
		}
	})

	t.Run("placed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

		svc := NewService(repo, logger)
		n := newOrderNumber(t)
		uID := newUserID(t)

		o, err := svc.Place(ctx, n, uID)
		if err != nil {
			t.Fatalf("could not place order: %s", err)
		}

		if o.Number != n {
			t.Errorf("could not match number")
			t.Errorf("got: %s", o.Number)
			t.Errorf("want: %s", n)
		}

		if o.Status != order.Placed {
			t.Errorf("could not match status")
			t.Errorf("got: %s", o.Status)
			t.Errorf("want: %s", order.Placed)
		}

		if o.PlacedBy != uID {
			t.Errorf("could not match placed by")
			t.Errorf("got: %s", o.PlacedBy)
			t.Errorf("want: %s", uID)
		}

		if time.Until(o.PlacedAt) > time.Second {
			t.Errorf("could not match placed at time: %s", o.PlacedAt)
		}
	})
}

func TestService_MarkAsShipped(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		svc := NewService(repo, logger)
		id := newID(t)

		repo.EXPECT().Find(ctx, id).Return(nil, order.ErrNotFound)

		o, err := svc.MarkAsShipped(ctx, id)
		if !errors.Is(err, ErrNotPlaced) && !errors.Is(err, order.ErrNotFound) {
			t.Fatalf("could match error: %s", err)
		}

		if o != nil {
			t.Fatalf("could not match a nil order: %v", o)
		}
	})

	t.Run("not added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		o := newPlacedOrder(t)

		repo.EXPECT().Find(ctx, o.ID).Return(o, nil)
		repo.EXPECT().Add(ctx, gomock.Any()).Return(order.ErrNotAdded)

		svc := NewService(repo, logger)

		o, err := svc.MarkAsShipped(ctx, o.ID)
		if !errors.Is(err, ErrNotPlaced) && !errors.Is(err, order.ErrNotAdded) {
			t.Fatalf("could match error: %s", err)
		}

		if o != nil {
			t.Fatalf("could not match a nil order: %v", o)
		}
	})

	t.Run("marked as shipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		o := newPlacedOrder(t)

		repo.EXPECT().Find(ctx, o.ID).Return(o, nil)
		repo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

		svc := NewService(repo, logger)

		o, err := svc.MarkAsShipped(ctx, o.ID)
		if err != nil {
			t.Fatalf("could match error: %s", err)
		}

		if o.Status != order.Shipped {
			t.Errorf("could not match status")
			t.Errorf("got: %s", o.Status)
			t.Errorf("want: %s", order.Shipped)
		}

		if time.Until(o.ShippedAt) > time.Second {
			t.Errorf("could not match shipped at time: %s", o.PlacedAt)
		}
	})
}

func TestService_MarkAsDelivered(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		svc := NewService(repo, logger)
		id := newID(t)

		repo.EXPECT().Find(ctx, id).Return(nil, order.ErrNotFound)

		o, err := svc.MarkAsDelivered(ctx, id)
		if !errors.Is(err, ErrNotMarkedAsDelivered) && !errors.Is(err, order.ErrNotFound) {
			t.Fatalf("could match error: %s", err)
		}

		if o != nil {
			t.Fatalf("could not match a nil order: %v", o)
		}
	})

	t.Run("not added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		o := newShippedOrder(t)

		repo.EXPECT().Find(ctx, o.ID).Return(o, nil)
		repo.EXPECT().Add(ctx, gomock.Any()).Return(order.ErrNotAdded)

		svc := NewService(repo, logger)

		o, err := svc.MarkAsDelivered(ctx, o.ID)
		if !errors.Is(err, ErrNotMarkedAsDelivered) && !errors.Is(err, order.ErrNotAdded) {
			t.Fatalf("could match error: %s", err)
		}

		if o != nil {
			t.Fatalf("could not match a nil order: %v", o)
		}
	})

	t.Run("marked as shipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		o := newShippedOrder(t)

		repo.EXPECT().Find(ctx, o.ID).Return(o, nil)
		repo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

		svc := NewService(repo, logger)

		o, err := svc.MarkAsDelivered(ctx, o.ID)
		if err != nil {
			t.Fatalf("could match error: %s", err)
		}

		if o.Status != order.Delivered {
			t.Errorf("could not match status")
			t.Errorf("got: %s", o.Status)
			t.Errorf("want: %s", order.Delivered)
		}

		if time.Until(o.DeliveredAt) > time.Second {
			t.Errorf("could not match delivered at time: %s", o.DeliveredAt)
		}
	})
}

func newPlacedOrder(t *testing.T) *order.Order {
	t.Helper()

	return order.Place(newOrderNumber(t), newUserID(t))
}

func newShippedOrder(t *testing.T) *order.Order {
	t.Helper()

	o := order.Place(newOrderNumber(t), newUserID(t))
	if err := o.MarkAsShipped(); err != nil {
		t.Fatalf("could not mark order as shipped: %s", err)
	}

	return o
}

func newID(t *testing.T) order.ID {
	t.Helper()

	return order.NewID()
}

func newUserID(t *testing.T) order.UserID {
	t.Helper()

	id, err := order.ParseUserID(uuid.NewString())
	if err != nil {
		t.Fatalf("could not generate user id: %s", err)
	}

	return id
}

func newOrderNumber(t *testing.T) order.Number {
	t.Helper()

	return order.GenerateNumber()
}
