// Package stock is the composition facade for the stock module.
// Its use cases depend on the product repository (to read and adjust product
// stock), which is injected by the application root.
package stock

import (
	"inventario/internal/modules/stock/application"
	domain "inventario/internal/modules/stock/domain"
	"inventario/internal/modules/stock/handler"
	"inventario/internal/modules/stock/infrastructure"

	productdomain "inventario/internal/modules/product/domain"

	"github.com/jmoiron/sqlx"
)

// Module holds the public surface of the stock module.
type Module struct {
	Handler *handler.StockHandler
	repo    domain.Repository
}

// New builds the stock module's object graph. productRepo is supplied by the
// application root so stock movements stay consistent with product stock.
func New(db *sqlx.DB, productRepo productdomain.Repository) *Module {
	repo := infrastructure.NewStockRepository(db)

	h := handler.NewStockHandler(
		application.NewAdjustUseCase(repo, productRepo),
		application.NewListMovementsUseCase(repo, productRepo),
	)

	return &Module{Handler: h, repo: repo}
}

// Repository exposes the stock repository for the sale module, which records
// stock movements when a sale is created.
func (m *Module) Repository() domain.Repository { return m.repo }
