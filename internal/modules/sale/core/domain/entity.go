package domain

import "time"

type Sale struct {
	ID        int64
	Total     int64
	Items     []SaleItem
	CreatedAt time.Time
}

// SaleProduct is the sale module's view of a sold product: just enough to
// identify it. Prices live on the item, since the product may change later.
type SaleProduct struct {
	ID   int64
	Name string
	Code string
}

type SaleItem struct {
	ID        int64
	SaleID    int64
	ProductID int64
	// Product is the related product. It is loaded when reading a sale and
	// nil otherwise; Save ignores it.
	Product   *SaleProduct
	Quantity  int
	UnitPrice int64
	Subtotal  int64
	CreatedAt time.Time
}
