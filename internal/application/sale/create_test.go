package sale_test

import (
	"errors"
	appSale "inventario/internal/application/sale"
	domainProduct "inventario/internal/domain/product"
	domainSale "inventario/internal/domain/sale"
	domainStock "inventario/internal/domain/stock"
	"testing"
)

func TestCreateUseCase_Execute(t *testing.T) {
	seedProduct := func(repo *mockProductRepo, stock int, active bool) int64 {
		p := &domainProduct.Product{
			Name:   "Laptop",
			Code:   "LAP-001",
			Price:  999.99,
			Stock:  stock,
			Active: active,
		}
		repo.seed(p)
		return p.ID
	}

	t.Run("crea venta con un item", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, true)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		result, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 2},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID == 0 {
			t.Error("expected sale ID to be set")
		}
		if len(result.Items) != 1 {
			t.Errorf("got %d items, want 1", len(result.Items))
		}
		expectedTotal := 999.99 * 2
		if result.Total != expectedTotal {
			t.Errorf("got total %.2f, want %.2f", result.Total, expectedTotal)
		}
	})

	t.Run("crea venta con múltiples items", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id1 := seedProduct(productRepo, 10, true)
		p2 := &domainProduct.Product{Name: "Mouse", Code: "MOU-001", Price: 29.99, Stock: 5, Active: true}
		productRepo.seed(p2)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		result, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id1, Quantity: 1},
				{ProductID: p2.ID, Quantity: 2},
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 {
			t.Errorf("got %d items, want 2", len(result.Items))
		}
		expectedTotal := 999.99 + (29.99 * 2)
		if result.Total != expectedTotal {
			t.Errorf("got total %.2f, want %.2f", result.Total, expectedTotal)
		}
	})

	t.Run("descuenta stock al crear venta", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, true)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 3},
			},
		})

		p, _ := productRepo.FindByID(id)
		if p.Stock != 7 {
			t.Errorf("got stock %d, want 7", p.Stock)
		}
	})

	t.Run("registra movimiento de stock tipo sale", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, true)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 2},
			},
		})

		if len(stockRepo.movements) != 1 {
			t.Fatalf("got %d movements, want 1", len(stockRepo.movements))
		}
		mv := stockRepo.movements[0]
		if mv.MovementType != domainStock.MovementSale {
			t.Errorf("got type %q, want %q", mv.MovementType, domainStock.MovementSale)
		}
		if mv.Quantity != -2 {
			t.Errorf("got quantity %d, want -2", mv.Quantity)
		}
	})

	t.Run("falla sin items", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{Items: []appSale.SaleItemInput{}})

		if !errors.Is(err, domainSale.ErrNoItems) {
			t.Errorf("got %v, want ErrNoItems", err)
		}
	})

	t.Run("falla con producto inexistente", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: 999, Quantity: 1},
			},
		})

		if !errors.Is(err, domainProduct.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("falla con stock insuficiente", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 2, true)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 5},
			},
		})

		if !errors.Is(err, domainProduct.ErrInsufficientStock) {
			t.Errorf("got %v, want ErrInsufficientStock", err)
		}
	})

	t.Run("falla con producto inactivo", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, false)
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 1},
			},
		})

		if err == nil {
			t.Fatal("expected error for inactive product, got nil")
		}
	})

	t.Run("propaga error del repositorio de ventas", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, true)
		saleRepo.saveErr = errors.New("db error")
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 1},
			},
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("propaga error al registrar movimiento de stock", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		productRepo := newMockProductRepo()
		stockRepo := newMockStockRepo()
		id := seedProduct(productRepo, 10, true)
		stockRepo.saveErr = errors.New("db error")
		uc := appSale.NewCreateUseCase(saleRepo, productRepo, stockRepo)

		_, err := uc.Execute(appSale.CreateInput{
			Items: []appSale.SaleItemInput{
				{ProductID: id, Quantity: 1},
			},
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
