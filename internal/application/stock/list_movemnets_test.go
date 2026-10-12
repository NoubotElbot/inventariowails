package stock_test

import (
	"errors"
	appStock "inventario/internal/application/stock"
	domainProduct "inventario/internal/domain/product"
	domain "inventario/internal/domain/stock"
	"testing"
)

func TestListMovementsUseCase_Execute(t *testing.T) {
	seedProduct := func(productRepo *mockProductRepo) int64 {
		p := &domainProduct.Product{
			Name:  "Laptop",
			Code:  "LAP-001",
			Price: 999,
			Stock: 10,
		}
		productRepo.seed(p)
		return p.ID
	}

	t.Run("retorna movimientos del producto", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo)
		stockRepo.Save(&domain.Movement{ProductID: id, Quantity: 5, MovementType: domain.MovementAdjustment})
		stockRepo.Save(&domain.Movement{ProductID: id, Quantity: -2, MovementType: domain.MovementSale})
		uc := appStock.NewListMovementsUseCase(stockRepo, productRepo)

		result, err := uc.Execute(id)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("got %d movements, want 2", len(result))
		}
	})

	t.Run("retorna lista vacía sin movimientos", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo)
		uc := appStock.NewListMovementsUseCase(stockRepo, productRepo)

		result, err := uc.Execute(id)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d movements, want 0", len(result))
		}
	})

	t.Run("no retorna movimientos de otro producto", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id1 := seedProduct(productRepo)
		id2 := seedProduct(productRepo)
		stockRepo.Save(&domain.Movement{ProductID: id2, Quantity: 5, MovementType: domain.MovementAdjustment})
		uc := appStock.NewListMovementsUseCase(stockRepo, productRepo)

		result, err := uc.Execute(id1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d movements for product %d, want 0", len(result), id1)
		}
	})

	t.Run("falla con producto inexistente", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		uc := appStock.NewListMovementsUseCase(stockRepo, productRepo)

		_, err := uc.Execute(999)

		if !errors.Is(err, domainProduct.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio de movimientos", func(t *testing.T) {
		stockRepo := newMockStockRepo()
		productRepo := newMockProductRepo()
		id := seedProduct(productRepo)
		stockRepo.findErr = errors.New("db error")
		uc := appStock.NewListMovementsUseCase(stockRepo, productRepo)

		_, err := uc.Execute(id)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
