package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestListUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewListUseCase(repo)

		result, err := uc.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d items, want 0", len(result))
		}
	})

	t.Run("retorna categorías existentes", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Category{Name: "Electrónica"})
		repo.Save(&domain.Category{Name: "Ropa"})
		uc := appCategory.NewListUseCase(repo)

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
		uc := appCategory.NewListUseCase(repo)

		_, err := uc.Execute()

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
