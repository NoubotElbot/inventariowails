package stock

func (m *Movement) Validate() error {
	if m.ProductID == 0 {
		return ErrProductRequired
	}
	if m.Quantity == 0 {
		return ErrQuantityZero
	}
	switch m.MovementType {
	case MovementAdjustment, MovementSale, MovementReturn:
		return nil
	default:
		return ErrInvalidMovement
	}
}
