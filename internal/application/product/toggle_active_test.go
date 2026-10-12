package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestToggleActiveUseCase_Execute(t *testing.T) {
	t.Run("desactiva producto activo", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Active: true}
		repo.Save(p)
		uc := appProduct.NewToggleActiveUseCase(repo)

		result, err := uc.Execute(p.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Active {
			t.Error("expected product to be inactive")
		}
	})

	t.Run("activa producto inactivo", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Active: false}
		repo.Save(p)
		uc := appProduct.NewToggleActiveUseCase(repo)

		result, err := uc.Execute(p.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Active {
			t.Error("expected product to be active")
		}
	})

	t.Run("doble toggle restaura estado original", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Active: true}
		repo.Save(p)
		uc := appProduct.NewToggleActiveUseCase(repo)

		uc.Execute(p.ID)                // desactiva
		result, err := uc.Execute(p.ID) // reactiva

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Active {
			t.Error("expected product to be active after double toggle")
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewToggleActiveUseCase(repo)

		_, err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Active: true}
		repo.Save(p)
		repo.updateErr = errors.New("db error")
		uc := appProduct.NewToggleActiveUseCase(repo)

		_, err := uc.Execute(p.ID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
