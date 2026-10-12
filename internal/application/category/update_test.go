package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	t.Run("actualiza categoría existente", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewUpdateUseCase(repo)

		result, err := uc.Execute(appCategory.UpdateInput{
			ID:   c.ID,
			Name: "Electrónica y más",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Name != "Electrónica y más" {
			t.Errorf("got %q, want %q", result.Name, "Electrónica y más")
		}
	})

	t.Run("falla con nombre vacío", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		uc := appCategory.NewUpdateUseCase(repo)

		_, err := uc.Execute(appCategory.UpdateInput{ID: c.ID, Name: ""})

		if !errors.Is(err, domain.ErrNameRequired) {
			t.Errorf("got %v, want ErrNameRequired", err)
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewUpdateUseCase(repo)

		_, err := uc.Execute(appCategory.UpdateInput{ID: 999, Name: "Algo"})

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio en update", func(t *testing.T) {
		repo := newMockRepo()
		c := &domain.Category{Name: "Electrónica"}
		repo.Save(c)
		repo.updateErr = errors.New("db error")
		uc := appCategory.NewUpdateUseCase(repo)

		_, err := uc.Execute(appCategory.UpdateInput{ID: c.ID, Name: "Nuevo nombre"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
