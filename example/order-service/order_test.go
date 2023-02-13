package order_test

import (
	"errors"
	"testing"
	"time"

	. "github.com/organization/order-service"

	"github.com/google/uuid"
)

func TestPlace(t *testing.T) {
	n := GenerateNumber()
	uID := userIDHelper(t)

	o := Place(n, uID)

	if o.ID.IsZero() {
		t.Errorf("could not match id: is zero: %s", o.ID)
	}

	if o.Number != n {
		t.Error("could not match number")
		t.Errorf("got: %s", n)
		t.Errorf("want: %s", o.Number)
	}

	if o.Status != Placed {
		t.Errorf("could not match placed status: %s", o.Status)
	}

	if o.PlacedBy != uID {
		t.Error("could not match placed by")
		t.Errorf("got: %s", uID)
		t.Errorf("want: %s", o.PlacedBy)
	}

	if time.Until(o.PlacedAt) > time.Second {
		t.Errorf("could not match placed at time: %s", o.PlacedAt)
	}

	if !o.ShippedAt.IsZero() {
		t.Errorf("could not match shipped at time as zero: %s", o.ShippedAt)
	}

	if !o.DeliveredAt.IsZero() {
		t.Errorf("could not match delivered at time as zero: %s", o.DeliveredAt)
	}
}

func TestOrder_MarkAsShipped(t *testing.T) {
	id := NewID()
	n := GenerateNumber()
	placedAt := time.Now()
	uID := userIDHelper(t)

	t.Run("placed", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Placed,
			PlacedBy:    uID,
			PlacedAt:    placedAt,
			ShippedAt:   time.Time{},
			DeliveredAt: time.Time{},
		}

		if err := o.MarkAsShipped(); err != nil {
			t.Fatalf("could not mark the order as shipped: %s", err)
		}

		if o.ID != id {
			t.Error("could not match id")
			t.Errorf("got: %s", id)
			t.Errorf("want: %s", o.ID)
		}

		if o.Number != n {
			t.Error("could not match number")
			t.Errorf("got: %s", n)
			t.Errorf("want: %s", o.Number)
		}

		if o.Status != Shipped {
			t.Errorf("could not match shipped status: %s", o.Status)
		}

		if o.PlacedBy != uID {
			t.Error("could not match original placed by")
			t.Errorf("got: %s", uID)
			t.Errorf("want: %s", o.PlacedBy)
		}

		if o.PlacedAt != placedAt {
			t.Error("could not match placed at")
			t.Errorf("got: %s", placedAt)
			t.Errorf("want: %s", o.PlacedBy)
		}

		if time.Until(o.ShippedAt) > time.Second {
			t.Errorf("could not match shipped at time: %s", o.ShippedAt)
		}

		if !o.DeliveredAt.IsZero() {
			t.Errorf("could not match delivered at time as zero: %s", o.DeliveredAt)
		}
	})

	t.Run("shipped", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Shipped,
			PlacedBy:    uID,
			PlacedAt:    time.Now(),
			ShippedAt:   time.Now(),
			DeliveredAt: time.Time{},
		}

		if err := o.MarkAsShipped(); !errors.Is(err, ErrNotShipped) {
			t.Fatalf("could mark the order as shipped: %s", err)
		}
	})

	t.Run("delivered", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Delivered,
			PlacedBy:    uID,
			PlacedAt:    placedAt,
			ShippedAt:   time.Now(),
			DeliveredAt: time.Now(),
		}

		if err := o.MarkAsShipped(); !errors.Is(err, ErrNotShipped) {
			t.Fatalf("could mark the order as shipped: %s", err)
		}
	})
}

func TestOrder_MarkAsDelivered(t *testing.T) {
	id := NewID()
	n := GenerateNumber()
	placedAt := time.Now()
	shippedAt := time.Now()
	uID := userIDHelper(t)

	t.Run("placed", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Placed,
			PlacedBy:    uID,
			PlacedAt:    placedAt,
			ShippedAt:   time.Time{},
			DeliveredAt: time.Time{},
		}

		if err := o.MarkAsDelivered(); !errors.Is(err, ErrNotDelivered) {
			t.Fatalf("could mark the order as delivered: %s", err)
		}
	})

	t.Run("shipped", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Shipped,
			PlacedBy:    uID,
			PlacedAt:    placedAt,
			ShippedAt:   shippedAt,
			DeliveredAt: time.Time{},
		}

		if err := o.MarkAsDelivered(); err != nil {
			t.Fatalf("could not mark the order as delivered: %s", err)
		}

		if o.ID != id {
			t.Error("could not match id")
			t.Errorf("got: %s", id)
			t.Errorf("want: %s", o.ID)
		}

		if o.Number != n {
			t.Error("could not match number")
			t.Errorf("got: %s", n)
			t.Errorf("want: %s", o.Number)
		}

		if o.Status != Delivered {
			t.Errorf("could not match delivered status: %s", o.Status)
		}

		if o.PlacedBy != uID {
			t.Error("could not match original placed by")
			t.Errorf("got: %s", uID)
			t.Errorf("want: %s", o.PlacedBy)
		}

		if o.PlacedAt != placedAt {
			t.Error("could not match placed at")
			t.Errorf("got: %s", placedAt)
			t.Errorf("want: %s", o.PlacedBy)
		}

		if o.ShippedAt != shippedAt {
			t.Error("could not match shipped at")
			t.Errorf("got: %s", shippedAt)
			t.Errorf("want: %s", o.ShippedAt)
		}

		if time.Until(o.DeliveredAt) > time.Second {
			t.Errorf("could not match delivered at time: %s", o.DeliveredAt)
		}
	})

	t.Run("delivered", func(t *testing.T) {
		o := &Order{
			ID:          id,
			Number:      n,
			Status:      Delivered,
			PlacedBy:    uID,
			PlacedAt:    placedAt,
			ShippedAt:   shippedAt,
			DeliveredAt: time.Now(),
		}

		if err := o.MarkAsShipped(); !errors.Is(err, ErrNotShipped) {
			t.Fatalf("could not mark the order as already delivered: %s", err)
		}
	})
}

func userIDHelper(t *testing.T) UserID {
	uID, err := ParseUserID(uuid.NewString())
	if err != nil {
		t.Fatalf("could ot parse user id: %s", err)
	}

	return uID
}
