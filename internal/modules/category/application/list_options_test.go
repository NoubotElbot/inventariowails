package application_test

import (
	"errors"
	appCategory "inventario/internal/modules/category/application"
	domain "inventario/internal/modules/category/domain"
	"testing"
)

func TestListOptionsUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		opts, err := appCategory.NewListOptionsUseCase(newMockRepo()).Execute()
		if err != nil || len(opts) != 0 {
			t.Fatalf("got %v, %v; want empty, nil", opts, err)
		}
	})

	t.Run("retorna id y nombre ordenados por nombre", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Category{Name: "Ropa"})
		repo.Save(&domain.Category{Name: "Electrónica"})

		opts, err := appCategory.NewListOptionsUseCase(repo).Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 2 || opts[0].Name != "Electrónica" || opts[1].Name != "Ropa" {
			t.Errorf("unexpected options: %v", opts)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		repo.findOptionsErr = errors.New("boom")
		if _, err := appCategory.NewListOptionsUseCase(repo).Execute(); !errors.Is(err, repo.findOptionsErr) {
			t.Errorf("got %v, want wrapped boom", err)
		}
	})
}
