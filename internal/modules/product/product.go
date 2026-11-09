// Package product is the composition facade for the product module.
// It wires the module's own layers and exposes both the Wails-facing handler
// and the product Repository consumed by the stock and sale modules.
package product

import (
	"inventario/internal/modules/product/application"
	domain "inventario/internal/modules/product/domain"
	"inventario/internal/modules/product/handler"
	"inventario/internal/modules/product/infrastructure"

	"github.com/jmoiron/sqlx"
)

// Module holds the public surface of the product module.
type Module struct {
	Handler *handler.ProductHandler
	repo    domain.Repository
}

// New builds the product module's object graph from a shared DB handle.
func New(db *sqlx.DB) *Module {
	repo := infrastructure.NewProductRepository(db)

	h := handler.NewProductHandler(
		application.NewListUseCase(repo),
		application.NewFindUseCase(repo),
		application.NewCreateUseCase(repo),
		application.NewUpdateUseCase(repo),
		application.NewDeleteUseCase(repo),
		application.NewToggleActiveUseCase(repo),
	)

	return &Module{Handler: h, repo: repo}
}

// Repository exposes the product repository for cross-module use cases
// (stock adjustments and sales) that need to read and adjust product stock.
func (m *Module) Repository() domain.Repository { return m.repo }
