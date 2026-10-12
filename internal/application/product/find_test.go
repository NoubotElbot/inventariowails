package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestFindUseCase_Execute(t *testing.T) {
	t.Run("encuentra producto existente", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999}
		repo.Save(p)
		uc := appProduct.NewFindUseCase(repo)

		result, err := uc.Execute(p.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Code != "LAP-001" {
			t.Errorf("got %q, want %q", result.Code, "LAP-001")
		}
	})

	t.Run("retorna ErrNotFound para ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewFindUseCase(repo)

		_, err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})
}
