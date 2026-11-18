package domain

import "time"

// MovementType says why stock moved and, with it, in which direction.
// Quantities are always positive; the type alone decides the sign.
type MovementType string

const (
	// MovementSale takes stock out (a sale).
	MovementSale MovementType = "sale"
	// MovementReturn puts stock back in (a customer return).
	MovementReturn MovementType = "return"
	// MovementAdjustmentIn is a manual count that found more stock than recorded.
	MovementAdjustmentIn MovementType = "adjustment_in"
	// MovementAdjustmentOut is a manual count that found less stock than recorded.
	MovementAdjustmentOut MovementType = "adjustment_out"
)

// IsManualAdjustment reports whether t can be registered by hand (stock
// counts). Sales and returns are registered by their own modules.
func (t MovementType) IsManualAdjustment() bool {
	return t == MovementAdjustmentIn || t == MovementAdjustmentOut
}

type Movement struct {
	ID           int64
	ProductID    int64
	Quantity     int // always > 0; direction comes from MovementType
	MovementType MovementType
	Notes        string
	CreatedAt    time.Time
}

// Delta is the movement's effect on product stock: +Quantity for inbound
// types (return, adjustment_in), -Quantity for outbound (sale,
// adjustment_out). It is only meaningful for a valid movement.
func (m *Movement) Delta() int {
	switch m.MovementType {
	case MovementSale, MovementAdjustmentOut:
		return -m.Quantity
	default:
		return m.Quantity
	}
}
