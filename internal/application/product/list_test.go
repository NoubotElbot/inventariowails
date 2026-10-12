package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestListUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d items, want 0", len(result))
		}
	})

	t.Run("retorna productos existentes", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Active: true})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", Price: 29, Active: true})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("got %d items, want 2", len(result))
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		repo.findAllErr = errors.New("db error")
		uc := appProduct.NewListUseCase(repo)

		_, err := uc.Execute()

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
