package application_test

import (
	"errors"
	appCategory "inventario/internal/modules/category/application"
	domain "inventario/internal/modules/category/domain"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	t.Run("elimina suavemente categoría existente", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewDeleteUseCase(repo)

		err := uc.Execute(c.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		stored, findErr := repo.FindByID(c.ID)
		if findErr != nil {
			t.Fatalf("unexpected error finding category: %v", findErr)
		}
		if stored.DeletedAt == nil {
			t.Error("expected category to be soft-deleted")
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

	t.Run("falla si la categoría ya está eliminada", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewDeleteUseCase(repo)
		uc.Execute(c.ID)

		err := uc.Execute(c.ID)

		if !errors.Is(err, domain.ErrAlreadyDeleted) {
			t.Errorf("got %v, want ErrAlreadyDeleted", err)
		}
	})

	t.Run("propaga error del repositorio en delete", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		repo.softDeleteErr = errors.New("db error")
		uc := appCategory.NewDeleteUseCase(repo)

		err := uc.Execute(c.ID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
