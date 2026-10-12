package product

import "time"

type Product struct {
	ID         int64
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      float64
	Stock      int
	MinStock   int
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Product) HasLowStock() bool {
	return p.Stock <= p.MinStock
}
