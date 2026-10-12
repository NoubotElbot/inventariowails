package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	t.Run("elimina producto existente", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999}
		repo.Save(p)
		uc := appProduct.NewDeleteUseCase(repo)

		err := uc.Execute(p.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err = repo.FindByID(p.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Error("expected product to be deleted")
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewDeleteUseCase(repo)

		err := uc.Execute(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999}
		repo.Save(p)
		repo.deleteErr = errors.New("db error")
		uc := appProduct.NewDeleteUseCase(repo)

		err := uc.Execute(p.ID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
