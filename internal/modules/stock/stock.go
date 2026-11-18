// Package stock is the composition root for the stock module.
// It wires the hexagonal layers: the GORM repository adapter and the product
// repository (injected by the application root) feed the core service through
// the outbound ports, and the service is injected into the Wails handler
// adapter through the inbound port. The handler is exposed to the
// application root, and the service to the sale module.
package stock

import (
	"inventario/internal/modules/stock/adapters/handler"
	"inventario/internal/modules/stock/adapters/repository"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/modules/stock/core/services"
	"inventario/internal/platform/gormdb"

	"gorm.io/gorm"
)

// Module holds the public surface of the stock module.
type Module struct {
	Handler *handler.StockHandler
	svc     ports.StockService
}

// New builds the stock module's object graph. products is the product
// module's repository; every stock change goes through the service, which
// updates the product's stock and records the movement in one unit of work.
func New(db *gorm.DB, products ports.ProductRepository) *Module {
	repo := repository.NewStockRepository(db)
	svc := services.NewStockService(gormdb.NewUnitOfWork(db), repo, products)

	return &Module{Handler: handler.NewStockHandler(svc), svc: svc}
}

// Service exposes the stock service so other modules (sale) can register
// movements. It joins the caller's unit of work when there is one.
func (m *Module) Service() ports.StockService { return m.svc }
