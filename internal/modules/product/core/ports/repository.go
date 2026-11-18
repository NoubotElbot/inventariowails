package ports

import (
	"context"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/platform/pagination"
)

// ProductRepository is the outbound (driven) port: what the core needs
// from persistence. Adapters translate storage errors into domain errors
// (e.g. domain.ErrNotFound, domain.ErrCodeDuplicate).
//
// It is also consumed by the stock and sale modules, which read products
// and adjust their stock.
type ProductRepository interface {
	FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Product], error)
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	FindByCode(ctx context.Context, code string) (*domain.Product, error)
	Save(ctx context.Context, p *domain.Product) error
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int64) error
	UpdateStock(ctx context.Context, id int64, stock int) error
}
