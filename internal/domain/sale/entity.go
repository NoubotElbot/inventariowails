package sale

import "time"

type Sale struct {
	ID        int64
	Total     float64
	Items     []SaleItem
	CreatedAt time.Time
}

type SaleItem struct {
	ID        int64
	SaleID    int64
	ProductID int64
	Quantity  int
	UnitPrice float64
	Subtotal  float64
	CreatedAt time.Time
}
