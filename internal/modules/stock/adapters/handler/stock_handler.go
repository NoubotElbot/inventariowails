// Package handler is the Wails driving adapter for the stock module.
// It translates frontend requests into calls on the ports.StockService
// inbound port and maps domain results into response DTOs.
package handler

import (
	"context"

	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/responses"
)

type StockHandler struct {
	svc ports.StockService
}

func NewStockHandler(svc ports.StockService) *StockHandler {
	return &StockHandler{svc: svc}
}

func (h *StockHandler) AdjustStock(ctx context.Context, req AdjustStockRequest) (*StockMovementResponse, error) {
	m, err := h.svc.AdjustStock(ctx, ports.AdjustInput{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Type:      domain.MovementType(req.MovementType),
		Notes:     req.Notes,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*m)
	return &r, nil
}

func (h *StockHandler) ListStockMovements(ctx context.Context, productID int64) ([]StockMovementResponse, error) {
	movements, err := h.svc.ListMovements(ctx, productID)
	if err != nil {
		return nil, err
	}

	return responses.MapSlice(movements, toResponse), nil
}
