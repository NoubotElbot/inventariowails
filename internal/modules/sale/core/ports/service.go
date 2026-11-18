package ports

import (
	"context"

	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/platform/pagination"
)

// SaleItemInput is one line of a sale: which product and how many units.
type SaleItemInput struct {
	ProductID int64
	Quantity  int
}

// CreateInput carries the data needed to create a sale.
type CreateInput struct {
	Items []SaleItemInput
}

// SaleService is the inbound (driving) port: the use cases the core exposes
// to adapters such as the Wails handler.
type SaleService interface {
	Create(ctx context.Context, input CreateInput) (*domain.Sale, error)
	List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Sale], error)
	Get(ctx context.Context, id int64) (*domain.Sale, error)
}
