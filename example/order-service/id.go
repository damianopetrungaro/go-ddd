package order

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrUserIDNotParsed represents an error returned by a user id value type (aka value objects)
	ErrUserIDNotParsed = errors.New("could not parse user id")
)

// ID represents an order id
type ID uuid.UUID

// NewID returns a new ID
func NewID() ID {
	return ID(uuid.New())
}

// ParseID returns a ID or an error if the given string is not a valid UserID
func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, ErrUserIDNotParsed
	}

	return ID(id), nil
}

// IsZero reports whether id represents the zero ID
func (id ID) IsZero() bool {
	return id == ID{}
}

// String returns a string representation of the ID
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// UserID represents the user id that submitted the order
type UserID uuid.UUID

// ParseUserID returns a UserID or an error if the given string is not a valid UserID
func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, ErrUserIDNotParsed
	}

	return UserID(id), nil
}

// IsZero reports whether id represents the zero UserID
func (id UserID) IsZero() bool {
	return id == UserID{}
}

// String returns a string representation of the UserID
func (id UserID) String() string {
	return uuid.UUID(id).String()
}
