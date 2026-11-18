package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
)

func validCreateInput() ports.CreateInput {
	return ports.CreateInput{
		Name:  "Laptop",
		Code:  "LAP-001",
		Price: 999,
		Stock: 10,
	}
}

func TestProductService_Create(t *testing.T) {
	t.Run("crea producto válido y activo", func(t *testing.T) {
		svc, repo := newTestService(t)

		got, err := svc.Create(t.Context(), validCreateInput())

		require.NoError(t, err)
		assert.NotZero(t, got.ID)
		assert.True(t, got.Active, "debe quedar activo por defecto")
		assert.False(t, got.CreatedAt.IsZero())
		assert.Contains(t, repo.products, got.ID)
	})

	t.Run("crea producto con categoría", func(t *testing.T) {
		svc, _ := newTestService(t)
		catID := int64(1)
		input := validCreateInput()
		input.CategoryID = &catID

		got, err := svc.Create(t.Context(), input)

		require.NoError(t, err)
		require.NotNil(t, got.CategoryID)
		assert.Equal(t, catID, *got.CategoryID)
	})

	t.Run("falla con código duplicado", func(t *testing.T) {
		svc, repo := newTestService(t)
		seedProduct(t, repo, "LAP-001")

		_, err := svc.Create(t.Context(), validCreateInput())

		require.ErrorIs(t, err, domain.ErrCodeDuplicate)
	})

	tests := []struct {
		name          string
		mutate        func(*ports.CreateInput)
		findByCodeErr error
		saveErr       error
		wantErr       error
	}{
		{name: "falla con nombre vacío", mutate: func(in *ports.CreateInput) { in.Name = "" }, wantErr: domain.ErrNameRequired},
		{name: "falla con código vacío", mutate: func(in *ports.CreateInput) { in.Code = "" }, wantErr: domain.ErrCodeRequired},
		{name: "falla con precio inválido", mutate: func(in *ports.CreateInput) { in.Price = 0 }, wantErr: domain.ErrPriceInvalid},
		{name: "falla con stock negativo", mutate: func(in *ports.CreateInput) { in.Stock = -1 }, wantErr: domain.ErrStockNegative},
		{name: "propaga error al verificar código", findByCodeErr: errDB, wantErr: errDB},
		{name: "propaga error del repositorio en save", saveErr: errDB, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			repo.findByCodeErr = tt.findByCodeErr
			repo.saveErr = tt.saveErr
			input := validCreateInput()
			if tt.mutate != nil {
				tt.mutate(&input)
			}

			_, err := svc.Create(t.Context(), input)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
