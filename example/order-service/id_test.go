package order_test

import (
	"errors"
	"testing"

	. "github.com/organization/order-service"

	"github.com/google/uuid"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if _, err := uuid.Parse(uuid.UUID(id).String()); err != nil {
		t.Fatalf("could not parse id")
	}
}

func TestID_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		if !(ID{}).IsZero() {
			t.Fatalf("could not match zero value id")

		}
	})

	t.Run("value", func(t *testing.T) {
		if NewID().IsZero() {
			t.Fatalf("could match zero value id")
		}
	})
}

func TestID_String(t *testing.T) {
	raw := uuid.New()
	id := ID(raw)
	if id.String() != raw.String() {
		t.Error("could not match id as string")
		t.Errorf("got: %s", id.String())
		t.Fatalf("want: %s", raw.String())
	}
}

func TestParseUserID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		raw := uuid.NewString()
		id, err := ParseUserID(raw)
		if err != nil {
			t.Fatalf("could not parse user id: %s", err)
		}

		if got := uuid.UUID(id).String(); got != raw {
			t.Error("could not match user id as its raw format")
			t.Errorf("got: %s", got)
			t.Fatalf("want: %s", raw)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		raw := "an invalid id"
		id, err := ParseUserID(raw)
		if !errors.Is(err, ErrUserIDNotParsed) {
			t.Fatalf("could not parse match error: %s", err)
		}
		if got := uuid.UUID(id).String(); got != (uuid.UUID{}).String() {
			t.Fatalf("could not match an empty user id: %s", got)
		}
	})
}

func TestUserID_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		if !(UserID{}).IsZero() {
			t.Fatalf("could not match zero value user id")
		}
	})

	t.Run("value", func(t *testing.T) {
		if UserID(uuid.New()).IsZero() {
			t.Fatalf("could match zero value user id")
		}
	})
}

func TestUserID_String(t *testing.T) {
	raw := uuid.New()
	id := UserID(raw)
	if id.String() != raw.String() {
		t.Error("could not match user id as string")
		t.Errorf("got: %s", raw.String())
		t.Fatalf("want: %s", id.String())
	}
}
