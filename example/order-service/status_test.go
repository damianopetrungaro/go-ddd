package order_test

import (
	"testing"

	. "github.com/organization/order-service"
)

func TestStatus_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		if !Status("").IsZero() {
			t.Fatalf("could not match zero value status")
		}
	})

	t.Run("value", func(t *testing.T) {
		if Shipped.IsZero() {
			t.Fatalf("could match zero value status")
		}
	})
}
