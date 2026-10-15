package sale

func (s *Sale) Validate() error {
	if len(s.Items) == 0 {
		return ErrNoItems
	}
	for _, item := range s.Items {
		if item.ProductID == 0 {
			return ErrProductInvalid
		}
		if item.Quantity <= 0 {
			return ErrQuantityInvalid
		}
	}
	return nil
}
