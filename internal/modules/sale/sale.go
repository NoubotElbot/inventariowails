// Package sale is the composition root for the sale module.
// It wires the hexagonal layers: the GORM repository adapter, the product
// repository and the stock service (both injected by the application root)
// feed the core service through the outbound ports, and the service is
// injected into the Wails handler adapter through the inbound port. Only the
// handler is exposed to the application root.
package sale

import (
	"inventario/internal/modules/sale/adapters/handler"
	"inventario/internal/modules/sale/adapters/repository"
	"inventario/internal/modules/sale/core/ports"
	"inventario/internal/modules/sale/core/services"
	"inventario/internal/platform/gormdb"

	"gorm.io/gorm"
)

// Module holds the public surface of the sale module.
type Module struct {
	Handler *handler.SaleHandler
}

// New builds the sale module's object graph. products and stock are supplied
// by the application root; creating a sale saves it and moves stock in one
// unit of work.
func New(db *gorm.DB, products ports.ProductRepository, stock ports.StockMovementRegistrar) *Module {
	repo := repository.NewSaleRepository(db)
	svc := services.NewSaleService(gormdb.NewUnitOfWork(db), repo, products, stock)

	return &Module{Handler: handler.NewSaleHandler(svc)}
}
