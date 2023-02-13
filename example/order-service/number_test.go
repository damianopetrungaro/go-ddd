package order_test

import (
	"testing"

	. "github.com/organization/order-service"
)

func TestGenerateNumber(t *testing.T) {
	n1 := GenerateNumber()
	if string(n1[:]) == "" {
		t.Fatalf("could match order number as an empty string")
	}

	n2 := GenerateNumber()
	if string(n2[:]) == "" {
		t.Fatalf("could match order number as an empty string")
	}

	if n1 == n2 {
		t.Fatalf("could match random order numbers")
	}
}

func TestNumber_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		if !(Number{}).IsZero() {
			t.Fatalf("could not match zero value number")
		}
	})

	t.Run("value", func(t *testing.T) {
		if GenerateNumber().IsZero() {
			t.Fatalf("could match zero value number")
		}
	})
}
