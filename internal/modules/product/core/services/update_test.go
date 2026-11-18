package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
)

func validUpdateInput(id int64) ports.UpdateInput {
	return ports.UpdateInput{
		ID:       id,
		Name:     "Laptop Pro",
		Code:     "LAP-002",
		Barcode:  "123456",
		Price:    1299,
		MinStock: 3,
	}
}

func TestProductService_Update(t *testing.T) {
	t.Run("actualiza producto existente sin tocar el stock", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")
		catID := int64(7)
		input := validUpdateInput(p.ID)
		input.CategoryID = &catID

		got, err := svc.Update(t.Context(), input)

		require.NoError(t, err)
		assert.Equal(t, "Laptop Pro", got.Name)
		assert.Equal(t, "LAP-002", got.Code)
		assert.Equal(t, "123456", got.Barcode)
		assert.Equal(t, int64(1299), got.Price)
		assert.Equal(t, 3, got.MinStock)
		assert.Equal(t, &catID, got.CategoryID)
		assert.Equal(t, p.Stock, got.Stock, "el stock solo cambia por movimientos")
		assert.Equal(t, *got, repo.products[p.ID])
	})

	t.Run("permite mantener el mismo código", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")
		input := validUpdateInput(p.ID)
		input.Code = p.Code

		_, err := svc.Update(t.Context(), input)

		require.NoError(t, err)
	})

	t.Run("falla con código duplicado de otro producto", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")
		seedProduct(t, repo, "MOU-001")
		input := validUpdateInput(p.ID)
		input.Code = "MOU-001"

		_, err := svc.Update(t.Context(), input)

		require.ErrorIs(t, err, domain.ErrCodeDuplicate)
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		svc, _ := newTestService(t)

		_, err := svc.Update(t.Context(), validUpdateInput(999))

		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	tests := []struct {
		name          string
		mutate        func(*ports.UpdateInput)
		findByIDErr   error
		findByCodeErr error
		updateErr     error
		wantErr       error
	}{
		{name: "falla con nombre vacío", mutate: func(in *ports.UpdateInput) { in.Name = "" }, wantErr: domain.ErrNameRequired},
		{name: "falla con precio inválido", mutate: func(in *ports.UpdateInput) { in.Price = -10 }, wantErr: domain.ErrPriceInvalid},
		{name: "propaga error al buscar", findByIDErr: errDB, wantErr: errDB},
		{name: "propaga error al verificar código", findByCodeErr: errDB, wantErr: errDB},
		{name: "propaga error del repositorio en update", updateErr: errDB, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			p := seedProduct(t, repo, "LAP-001")
			repo.findByIDErr = tt.findByIDErr
			repo.findByCodeErr = tt.findByCodeErr
			repo.updateErr = tt.updateErr
			input := validUpdateInput(p.ID)
			if tt.mutate != nil {
				tt.mutate(&input)
			}

			_, err := svc.Update(t.Context(), input)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
