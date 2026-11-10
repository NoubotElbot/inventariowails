// Package sale is the composition facade for the sale module.
// Creating a sale orchestrates product (read price/stock, adjust stock) and
// stock (record movement), whose repositories are injected by the app root.
package sale

import (
	"inventario/internal/modules/sale/application"
	"inventario/internal/modules/sale/handler"
	"inventario/internal/modules/sale/infrastructure"

	productdomain "inventario/internal/modules/product/domain"
	stockdomain "inventario/internal/modules/stock/domain"

	"github.com/jmoiron/sqlx"
)

// Module holds the public surface of the sale module.
type Module struct {
	Handler *handler.SaleHandler
}

// New builds the sale module's object graph. productRepo and stockRepo are
// supplied by the application root.
func New(db *sqlx.DB, productRepo productdomain.Repository, stockRepo stockdomain.Repository) *Module {
	repo := infrastructure.NewSaleRepository(db)

	h := handler.NewSaleHandler(
		application.NewCreateUseCase(repo, productRepo, stockRepo),
		application.NewListUseCase(repo),
		application.NewFindUseCase(repo),
	)

	return &Module{Handler: h}
}
