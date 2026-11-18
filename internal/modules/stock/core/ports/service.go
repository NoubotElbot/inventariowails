package ports

import (
	"context"

	"inventario/internal/modules/stock/core/domain"
)

// RegisterMovementInput describes a stock movement. Quantity is always
// positive; Type decides the direction (sale and adjustment_out take stock
// out, return and adjustment_in put it in).
type RegisterMovementInput struct {
	ProductID int64
	Quantity  int
	Type      domain.MovementType
	Notes     string
}

// AdjustInput describes a manual stock adjustment. Type must be
// adjustment_in or adjustment_out.
type AdjustInput struct {
	ProductID int64
	Quantity  int
	Type      domain.MovementType
	Notes     string
}

// StockService is the inbound (driving) port: the use cases the core exposes
// to the Wails handler and to other modules (sale registers its movements
// through RegisterMovement).
type StockService interface {
	AdjustStock(ctx context.Context, input AdjustInput) (*domain.Movement, error)
	RegisterMovement(ctx context.Context, input RegisterMovementInput) (*domain.Movement, error)
	ListMovements(ctx context.Context, productID int64) ([]domain.Movement, error)
}
