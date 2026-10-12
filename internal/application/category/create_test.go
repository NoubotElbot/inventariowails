package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"testing"
)

func TestCreateUseCase_Execute(t *testing.T) {
	t.Run("crea categoría válida", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewCreateUseCase(repo)

		result, err := uc.Execute(appCategory.CreateInput{Name: "Electrónica"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID == 0 {
			t.Error("expected ID to be set after save")
		}
		if result.Name != "Electrónica" {
			t.Errorf("got %q, want %q", result.Name, "Electrónica")
		}
	})

	t.Run("falla con nombre vacío", func(t *testing.T) {
		repo := newMockRepo()
		uc := appCategory.NewCreateUseCase(repo)

		_, err := uc.Execute(appCategory.CreateInput{Name: ""})

		if !errors.Is(err, domain.ErrNameRequired) {
			t.Errorf("got %v, want ErrNameRequired", err)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		repo.saveErr = errors.New("db error")
		uc := appCategory.NewCreateUseCase(repo)

		_, err := uc.Execute(appCategory.CreateInput{Name: "Electrónica"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
