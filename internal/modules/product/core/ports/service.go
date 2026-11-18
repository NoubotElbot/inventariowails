package ports

import (
	"context"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/platform/pagination"
)

// CreateInput carries the data needed to create a product.
type CreateInput struct {
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      int64
	Stock      int
	MinStock   int
}

// UpdateInput carries the data needed to update a product. Stock is not
// editable here: it only changes through stock movements.
type UpdateInput struct {
	ID         int64
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      int64
	MinStock   int
}

// ProductService is the inbound (driving) port: the use cases the core
// exposes to adapters such as the Wails handler.
type ProductService interface {
	List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Product], error)
	Get(ctx context.Context, id int64) (*domain.Product, error)
	Create(ctx context.Context, input CreateInput) (*domain.Product, error)
	Update(ctx context.Context, input UpdateInput) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
	ToggleActive(ctx context.Context, id int64) (*domain.Product, error)
}
