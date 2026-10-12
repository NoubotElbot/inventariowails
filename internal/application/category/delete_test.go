package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	t.Run("elimina categoría existente", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewDeleteUseCase(repo)

		err := uc.Execute(c.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = repo.FindByID(c.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Error("expected category to be deleted")
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewDeleteUseCase(repo)

		err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio en delete", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		repo.deleteErr = errors.New("db error")
		uc := appCategory.NewDeleteUseCase(repo)

		err := uc.Execute(c.ID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
