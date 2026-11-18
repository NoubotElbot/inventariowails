package ports

import (
	"context"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/platform/pagination"
)

// CategoryRepository is the outbound (driven) port: what the core needs
// from persistence. Adapters translate storage errors into domain errors
// (e.g. domain.ErrNotFound, domain.ErrNameDuplicate).
type CategoryRepository interface {
	FindOptions(ctx context.Context) ([]domain.Option, error)
	FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Category], error)
	FindByID(ctx context.Context, id int64) (*domain.Category, error)
	Save(ctx context.Context, c *domain.Category) error
	Update(ctx context.Context, c *domain.Category) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
}
