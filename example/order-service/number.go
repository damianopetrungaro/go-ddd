package order

import (
	"strings"

	"github.com/google/uuid"
)

// Number represents an order number.
// It is a random 32 chars string
type Number [32]byte

// GenerateNumber generates a new random order Number
// It is a factory function that uses the ubiquitous language of the domain
func GenerateNumber() Number {
	var n [32]byte
	for i, c := range strings.ReplaceAll(uuid.NewString(), `-`, ``) {
		n[i] = byte(c)
	}
	return n
}

// IsZero reports whether n represents the zero Number
func (n Number) IsZero() bool {
	return n == [32]byte{}
}

// String returns the Number as string
func (n Number) String() string {
	return string(n[:])
}
