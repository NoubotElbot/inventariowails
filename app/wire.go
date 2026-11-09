package app

import (
	"inventario/internal/modules/category"
	"inventario/internal/modules/product"
	"inventario/internal/modules/sale"
	"inventario/internal/modules/stock"
)

// wire assembles the business modules and injects the cross-module
// dependencies (product stock feeds both stock adjustments and sales).
func (a *App) wire() {
	categoryMod := category.New(a.db)
	productMod := product.New(a.db)
	stockMod := stock.New(a.db, productMod.Repository())
	saleMod := sale.New(a.db, productMod.Repository(), stockMod.Repository())

	a.CategoryHandler = categoryMod.Handler
	a.ProductHandler = productMod.Handler
	a.StockHandler = stockMod.Handler
	a.SaleHandler = saleMod.Handler
}
