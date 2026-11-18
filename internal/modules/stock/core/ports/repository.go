package ports

import (
	"context"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/stock/core/domain"
)

// StockRepository is the outbound (driven) port for persisting stock
// movements. Implementations join the transaction carried by ctx, if any.
type StockRepository interface {
	Save(ctx context.Context, m *domain.Movement) error
	FindByProductID(ctx context.Context, productID int64) ([]domain.Movement, error)
}

// ProductRepository is the narrow outbound port for what stock needs from
// products: read the current stock and overwrite it. The product module's
// repository satisfies it. Implementations return productdomain.ErrNotFound
// for a missing product.
type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*productdomain.Product, error)
	UpdateStock(ctx context.Context, id int64, stock int) error
}
