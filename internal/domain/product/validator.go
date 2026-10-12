package product

import "strings"

func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(p.Code) == "" {
		return ErrCodeRequired
	}
	if p.Price <= 0 {
		return ErrPriceInvalid
	}
	if p.Stock < 0 {
		return ErrStockNegative
	}
	return nil
}
