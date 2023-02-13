package order

const (
	Placed    Status = "placed"
	Shipped   Status = "shipped"
	Delivered Status = "delivered"
)

// Status represent an order status
type Status string

// IsZero reports whether o represents the zero Status
func (s Status) IsZero() bool {
	return s == ""
}

// String returns the Status as string
func (s Status) String() string {
	return string(s)
}
