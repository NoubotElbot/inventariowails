package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"testing"
)

func TestCreateUseCase_Execute(t *testing.T) {
	validInput := func() appProduct.CreateInput {
		return appProduct.CreateInput{
			Name:  "Laptop",
			Code:  "LAP-001",
			Price: 999.99,
			Stock: 10,
		}
	}

	t.Run("crea producto válido", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)

		result, err := uc.Execute(validInput())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID == 0 {
			t.Error("expected ID to be set after save")
		}
		if !result.Active {
			t.Error("expected product to be active by default")
		}
	})

	t.Run("crea producto con categoría", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)
		catID := int64(1)
		input := validInput()
		input.CategoryID = &catID

		result, err := uc.Execute(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.CategoryID == nil || *result.CategoryID != catID {
			t.Errorf("expected CategoryID %d, got %v", catID, result.CategoryID)
		}
	})

	t.Run("falla con nombre vacío", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)
		input := validInput()
		input.Name = ""

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrNameRequired) {
			t.Errorf("got %v, want ErrNameRequired", err)
		}
	})

	t.Run("falla con código vacío", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)
		input := validInput()
		input.Code = ""

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrCodeRequired) {
			t.Errorf("got %v, want ErrCodeRequired", err)
		}
	})

	t.Run("falla con precio inválido", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)
		input := validInput()
		input.Price = 0

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrPriceInvalid) {
			t.Errorf("got %v, want ErrPriceInvalid", err)
		}
	})

	t.Run("falla con stock negativo", func(t *testing.T) {
		repo := newMockRepo()
		uc := appProduct.NewCreateUseCase(repo)
		input := validInput()
		input.Stock = -1

		_, err := uc.Execute(input)

		if !errors.Is(err, domain.ErrStockNegative) {
			t.Errorf("got %v, want ErrStockNegative", err)
		}
	})

	t.Run("falla con código duplicado", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999})
		uc := appProduct.NewCreateUseCase(repo)

		_, err := uc.Execute(validInput())

		if !errors.Is(err, domain.ErrCodeDuplicate) {
			t.Errorf("got %v, want ErrCodeDuplicate", err)
		}
	})

	t.Run("propaga error del repositorio en save", func(t *testing.T) {
		repo := newMockRepo()
		repo.saveErr = errors.New("db error")
		uc := appProduct.NewCreateUseCase(repo)

		_, err := uc.Execute(validInput())

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
