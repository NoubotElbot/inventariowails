package ports

import (
	"context"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/platform/pagination"
)

// CreateInput carries the data needed to create a category.
type CreateInput struct {
	Name string
}

// UpdateInput carries the data needed to update a category.
type UpdateInput struct {
	ID   int64
	Name string
}

// CategoryService is the inbound (driving) port: the use cases the core
// exposes to adapters such as the Wails handler.
type CategoryService interface {
	List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Category], error)
	ListOptions(ctx context.Context) ([]domain.Option, error)
	Get(ctx context.Context, id int64) (*domain.Category, error)
	Create(ctx context.Context, input CreateInput) (*domain.Category, error)
	Update(ctx context.Context, input UpdateInput) (*domain.Category, error)
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
}
