package sale_test

import (
	"errors"
	appSale "inventario/internal/application/sale"
	domainSale "inventario/internal/domain/sale"
	"testing"
)

func TestFindUseCase_Execute(t *testing.T) {
	t.Run("encuentra venta existente con items", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		s := &domainSale.Sale{
			Total: 999.99,
			Items: []domainSale.SaleItem{
				{ProductID: 1, Quantity: 1, UnitPrice: 999.99, Subtotal: 999.99},
			},
		}
		saleRepo.Save(s)
		uc := appSale.NewFindUseCase(saleRepo)

		result, err := uc.Execute(s.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 999.99 {
			t.Errorf("got total %.2f, want 999.99", result.Total)
		}
		if len(result.Items) != 1 {
			t.Errorf("got %d items, want 1", len(result.Items))
		}
	})

	t.Run("retorna ErrNotFound para ID inexistente", func(t *testing.T) {
		saleRepo := newMockSaleRepo()
		uc := appSale.NewFindUseCase(saleRepo)

		_, err := uc.Execute(999)

		if !errors.Is(err, domainSale.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})
}
