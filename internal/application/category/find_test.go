package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestFindUseCase_Execute(t *testing.T) {
	t.Run("encuentra categoría existente", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewFindUseCase(repo)

		result, err := uc.Execute(c.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Name != "Electrónica" {
			t.Errorf("got %q, want %q", result.Name, "Electrónica")
		}
	})

	t.Run("retorna ErrNotFound para ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewFindUseCase(repo)

		_, err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})
}
