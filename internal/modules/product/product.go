// Package product is the composition root for the product module.
// It wires the hexagonal layers: the GORM repository adapter (driven) is
// injected into the core service through the outbound port, and the service
// is injected into the Wails handler adapter (driving) through the inbound
// port. The handler is exposed to the application root, and the repository
// port to the stock and sale modules.
package product

import (
	"inventario/internal/modules/product/adapters/handler"
	"inventario/internal/modules/product/adapters/repository"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/modules/product/core/services"

	"gorm.io/gorm"
)

// Module holds the public surface of the product module.
type Module struct {
	Handler *handler.ProductHandler
	repo    ports.ProductRepository
}

// New builds the product module's object graph from a shared DB handle.
func New(db *gorm.DB) *Module {
	repo := repository.NewProductRepository(db)
	svc := services.NewProductService(repo)

	return &Module{Handler: handler.NewProductHandler(svc), repo: repo}
}

// Repository exposes the product repository for cross-module use cases
// (stock adjustments and sales) that need to read and adjust product stock.
func (m *Module) Repository() ports.ProductRepository { return m.repo }
