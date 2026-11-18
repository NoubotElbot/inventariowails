package handler

// AdjustStockRequest is a manual stock adjustment. Quantity is always
// positive; MovementType must be "adjustment_in" or "adjustment_out" (the
// service enforces both).
type AdjustStockRequest struct {
	ProductID    int64  `json:"product_id"`
	Quantity     int    `json:"quantity"`
	MovementType string `json:"movement_type"`
	Notes        string `json:"notes"`
}
