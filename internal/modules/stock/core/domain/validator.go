package domain

func (m *Movement) Validate() error {
	if m.ProductID == 0 {
		return ErrProductRequired
	}
	if m.Quantity <= 0 {
		return ErrQuantityNotPositive
	}
	switch m.MovementType {
	case MovementSale, MovementReturn, MovementAdjustmentIn, MovementAdjustmentOut:
		return nil
	default:
		return ErrInvalidMovement
	}
}
