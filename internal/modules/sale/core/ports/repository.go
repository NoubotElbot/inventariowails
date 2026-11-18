package ports

import (
	"context"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/sale/core/domain"
	stockdomain "inventario/internal/modules/stock/core/domain"
	stockports "inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/pagination"
)

// SaleRepository is the outbound (driven) port for persisting sales and
// their items. Adapters translate storage errors into domain errors (e.g.
// domain.ErrNotFound) and join the transaction carried by ctx, if any.
type SaleRepository interface {
	Save(ctx context.Context, s *domain.Sale) error
	// FindAll returns a page of sales, newest first, with their items.
	FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Sale], error)
	FindByID(ctx context.Context, id int64) (*domain.Sale, error)
}

// ProductRepository is the narrow outbound port for what sale needs from
// products: price, active flag and stock. The product module's repository
// satisfies it. Implementations return productdomain.ErrNotFound for a
// missing product.
type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*productdomain.Product, error)
}

// StockMovementRegistrar applies a stock movement to a product and records
// it. The stock module's service (stockports.StockService) satisfies it and
// joins the caller's unit of work.
type StockMovementRegistrar interface {
	RegisterMovement(ctx context.Context, in stockports.RegisterMovementInput) (*stockdomain.Movement, error)
}
