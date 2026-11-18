package handler

import (
	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/platform/responses"
)

type StockMovementResponse struct {
	ID           int64  `json:"id"`
	ProductID    int64  `json:"product_id"`
	Quantity     int    `json:"quantity"`
	MovementType string `json:"movement_type"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

func toResponse(m domain.Movement) StockMovementResponse {
	return StockMovementResponse{
		ID:           m.ID,
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		MovementType: string(m.MovementType),
		Notes:        m.Notes,
		CreatedAt:    responses.FormatTime(m.CreatedAt),
	}
}
