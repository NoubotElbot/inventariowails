package app

import (
	appCategory "inventario/internal/application/category"
	appProduct "inventario/internal/application/product"
	appStock "inventario/internal/application/stock"
	"inventario/internal/handler"
	"inventario/internal/infrastructure/sqlite"
)

func (a *App) wire() {
	// --- Repositorios ---
	categoryRepo := sqlite.NewCategoryRepository(a.db)
	productRepo := sqlite.NewProductRepository(a.db)
	stockRepo := sqlite.NewStockRepository(a.db)

	// --- Category ---
	a.CategoryHandler = handler.NewCategoryHandler(
		appCategory.NewListUseCase(categoryRepo),
		appCategory.NewFindUseCase(categoryRepo),
		appCategory.NewCreateUseCase(categoryRepo),
		appCategory.NewUpdateUseCase(categoryRepo),
		appCategory.NewDeleteUseCase(categoryRepo),
	)

	// --- Product ---
	a.ProductHandler = handler.NewProductHandler(
		appProduct.NewListUseCase(productRepo),
		appProduct.NewFindUseCase(productRepo),
		appProduct.NewCreateUseCase(productRepo),
		appProduct.NewUpdateUseCase(productRepo),
		appProduct.NewDeleteUseCase(productRepo),
		appProduct.NewToggleActiveUseCase(productRepo),
	)

	// --- Stock ---
	a.StockHandler = handler.NewStockHandler(
		appStock.NewAdjustUseCase(stockRepo, productRepo),
		appStock.NewListMovementsUseCase(stockRepo, productRepo),
	)
}
