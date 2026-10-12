package stock_test

import (
	"errors"
	appStock "inventario/internal/application/stock"
	domainProduct "inventario/internal/domain/product"
	"testing"
)

func TestAdjustUseCase_Execute(t *testing.T) {
	seedProduct := func(productRepo *mockProductRepo, stock int) int64 {
		p := &domainProduct.Product{
			Name:  "Laptop",
			Code:  "LAP-001",
			Price: 999,
			Stock: stock,
		}
		productRepo.seed(p)
		return p.ID
	}

	t.Run("ajuste positivo incrementa stock", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  5,
			Notes:     "reposición",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		p, _ := productRepo.FindByID(id)
		if p.Stock != 15 {
			t.Errorf("got stock %d, want 15", p.Stock)
		}
	})

	t.Run("ajuste negativo decrementa stock", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  -3,
			Notes:     "merma",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		p, _ := productRepo.FindByID(id)
		if p.Stock != 7 {
			t.Errorf("got stock %d, want 7", p.Stock)
		}
	})

	t.Run("registra movimiento en repositorio", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		m, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  5,
			Notes:     "reposición",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.ID == 0 {
			t.Error("expected movement ID to be set")
		}
		if len(stockRepo.movements) != 1 {
			t.Errorf("got %d movements, want 1", len(stockRepo.movements))
		}
	})

	t.Run("falla si stock resultante es negativo", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 5)
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  -10,
		})

		if !errors.Is(err, domainProduct.ErrInsufficientStock) {
			t.Errorf("got %v, want ErrInsufficientStock", err)
		}
	})

	t.Run("falla con cantidad cero", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  0,
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("falla con producto inexistente", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: 999,
			Quantity:  5,
		})

		if !errors.Is(err, domainProduct.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio de stock", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		stockRepo.saveErr = errors.New("db error")
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  5,
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("propaga error al actualizar stock del producto", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo, 10)
		productRepo.updateStockErr = errors.New("db error")
		uc := appStock.NewAdjustUseCase(stockRepo, productRepo)

		_, err := uc.Execute(appStock.AdjustInput{
			ProductID: id,
			Quantity:  5,
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
