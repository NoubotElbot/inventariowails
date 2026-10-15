package sale_test

import (
	"errors"
	appSale "inventario/internal/application/sale"
	domainSale "inventario/internal/domain/sale"
	"testing"
)

func TestListUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		uc := appSale.NewListUseCase(saleRepo)

		result, err := uc.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d items, want 0", len(result))
		}
	})

	t.Run("retorna ventas existentes", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		saleRepo.Save(&domainSale.Sale{
			Total: 100,
			Items: []domainSale.SaleItem{{ProductID: 1, Quantity: 1, UnitPrice: 100, Subtotal: 100}},
		})
		saleRepo.Save(&domainSale.Sale{
			Total: 200,
			Items: []domainSale.SaleItem{{ProductID: 2, Quantity: 2, UnitPrice: 100, Subtotal: 200}},
		})
		uc := appSale.NewListUseCase(saleRepo)

		result, err := uc.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("got %d items, want 2", len(result))
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		saleRepo.findErr = errors.New("db error")
		uc := appSale.NewListUseCase(saleRepo)

		_, err := uc.Execute()

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
