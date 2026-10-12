package handler

import (
	"errors"
	appStock "inventario/internal/application/stock"
	domainProduct "inventario/internal/domain/product"
	domain "inventario/internal/domain/stock"
)

// --- DTOs ---

type StockMovementResponse struct {
	ID           int64  `json:"id"`
	ProductID    int64  `json:"product_id"`
	Quantity     int    `json:"quantity"`
	MovementType string `json:"movement_type"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

type AdjustStockRequest struct {
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Notes     string `json:"notes"`
}

// --- Handler ---

type StockHandler struct {
	adjust        *appStock.AdjustUseCase
	listMovements *appStock.ListMovementsUseCase
}

func NewStockHandler(
	adjust *appStock.AdjustUseCase,
	listMovements *appStock.ListMovementsUseCase,
) *StockHandler {
	return &StockHandler{
		adjust:        adjust,
		listMovements: listMovements,
	}
}

// --- métodos expuestos a Wails ---

func (h *StockHandler) AdjustStock(req AdjustStockRequest) (*StockMovementResponse, error) {
	m, err := h.adjust.Execute(appStock.AdjustInput{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Notes:     req.Notes,
	})
	if err != nil {
		return nil, mapStockError(err)
	}

	r := toStockMovementResponse(*m)
	return &r, nil
}

func (h *StockHandler) ListStockMovements(productID int64) ([]StockMovementResponse, error) {
	movements, err := h.listMovements.Execute(productID)
	if err != nil {
		return nil, mapStockError(err)
	}

	response := make([]StockMovementResponse, 0, len(movements))
	for _, m := range movements {
		response = append(response, toStockMovementResponse(m))
	}

	return response, nil
}

// --- helpers privados ---

func toStockMovementResponse(m domain.Movement) StockMovementResponse {
	return StockMovementResponse{
		ID:           m.ID,
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		MovementType: string(m.MovementType),
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func mapStockError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domainProduct.ErrNotFound):
		return errors.New("producto no encontrado")
	case errors.Is(err, domainProduct.ErrInsufficientStock):
		return errors.New("stock insuficiente para el ajuste")
	case errors.Is(err, domain.ErrProductRequired):
		return errors.New("el producto es requerido")
	case errors.Is(err, domain.ErrQuantityZero):
		return errors.New("la cantidad no puede ser cero")
	case errors.Is(err, domain.ErrInvalidMovement):
		return errors.New("tipo de movimiento inválido")
	default:
		return errors.New("error interno del servidor")
	}
}
