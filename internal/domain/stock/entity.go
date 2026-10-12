package stock

import "time"

type MovementType string

const (
	MovementAdjustment MovementType = "adjustment"
	MovementSale       MovementType = "sale"
	MovementReturn     MovementType = "return"
)

type Movement struct {
	ID           int64
	ProductID    int64
	Quantity     int
	MovementType MovementType
	Notes        string
	CreatedAt    time.Time
}
