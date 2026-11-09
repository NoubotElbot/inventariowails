package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestRestoreUseCase_Execute(t *testing.T) {
	t.Run("restaura categoría eliminada", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		repo.SoftDelete(c.ID)
		uc := appCategory.NewRestoreUseCase(repo)

		err := uc.Execute(c.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		stored, _ := repo.FindByID(c.ID)
		if stored.DeletedAt != nil {
			t.Error("expected DeletedAt to be nil after restore")
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewRestoreUseCase(repo)

		err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("falla si la categoría no está eliminada", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewRestoreUseCase(repo)

		err := uc.Execute(c.ID)

		if !errors.Is(err, domain.ErrNotDeleted) {
			t.Errorf("got %v, want ErrNotDeleted", err)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		repo.SoftDelete(c.ID)
		repo.restoreErr = errors.New("db error")
		uc := appCategory.NewRestoreUseCase(repo)

		err := uc.Execute(c.ID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
