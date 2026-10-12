package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	seedProduct := func(repo *mockRepo) *domain.Product {
		p := &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999, Stock: 10}
		repo.Save(p)
		return p
	}

	validInput := func(id int64) appProduct.UpdateInput {
		return appProduct.UpdateInput{
			ID:    id,
			Name:  "Laptop Pro",
			Code:  "LAP-002",
			Price: 1299.99,
		}
	}

	t.Run("actualiza producto existente", func(t *testing.T) {
		repo := newMockRepo()
		p := seedProduct(repo)
		uc := appProduct.NewUpdateUseCase(repo)

		result, err := uc.Execute(validInput(p.ID))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Name != "Laptop Pro" {
			t.Errorf("got %q, want %q", result.Name, "Laptop Pro")
		}
		if result.Code != "LAP-002" {
			t.Errorf("got %q, want %q", result.Code, "LAP-002")
		}
	})

	t.Run("permite mantener el mismo código", func(t *testing.T) {
		repo := newMockRepo()
		p := seedProduct(repo)
		uc := appProduct.NewUpdateUseCase(repo)
		input := validInput(p.ID)
		input.Code = p.Code // mismo código, no debe fallar

		_, err := uc.Execute(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewUpdateUseCase(repo)

		_, err := uc.Execute(validInput(999))

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("falla con código duplicado de otro producto", func(t *testing.T) {
		repo := newMockRepo()
		p := seedProduct(repo)
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", Price: 29, Stock: 5})
		uc := appProduct.NewUpdateUseCase(repo)
		input := validInput(p.ID)
		input.Code = "MOU-001"

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrCodeDuplicate) {
			t.Errorf("got %v, want ErrCodeDuplicate", err)
		}
	})

	t.Run("falla con precio inválido", func(t *testing.T) {
		repo := newMockRepo()
		p := seedProduct(repo)
		uc := appProduct.NewUpdateUseCase(repo)
		input := validInput(p.ID)
		input.Price = -10

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrPriceInvalid) {
			t.Errorf("got %v, want ErrPriceInvalid", err)
		}
	})

	t.Run("propaga error del repositorio en update", func(t *testing.T) {
		repo := newMockRepo()
		p := seedProduct(repo)
		repo.updateErr = errors.New("db error")
		uc := appProduct.NewUpdateUseCase(repo)

		_, err := uc.Execute(validInput(p.ID))

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
