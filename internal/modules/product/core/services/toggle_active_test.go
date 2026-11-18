package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/product/core/domain"
)

func TestProductService_ToggleActive(t *testing.T) {
	t.Run("invierte el estado y lo persiste", func(t *testing.T) {
		for _, initial := range []bool{true, false} {
			svc, repo := newTestService(t)
			p := seedProduct(t, repo, "LAP-001")
			p.Active = initial
			repo.products[p.ID] = p

			got, err := svc.ToggleActive(t.Context(), p.ID)

			require.NoError(t, err)
			assert.Equal(t, !initial, got.Active)
			assert.Equal(t, !initial, repo.products[p.ID].Active)
		}
	})

	t.Run("doble toggle restaura estado original", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")

		_, err := svc.ToggleActive(t.Context(), p.ID)
		require.NoError(t, err)
		got, err := svc.ToggleActive(t.Context(), p.ID)

		require.NoError(t, err)
		assert.True(t, got.Active)
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		svc, _ := newTestService(t)

		_, err := svc.ToggleActive(t.Context(), 999)

		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	tests := []struct {
		name        string
		findByIDErr error
		updateErr   error
	}{
		{name: "propaga error al buscar", findByIDErr: errDB},
		{name: "propaga error del repositorio en update", updateErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			p := seedProduct(t, repo, "LAP-001")
			repo.findByIDErr = tt.findByIDErr
			repo.updateErr = tt.updateErr

			_, err := svc.ToggleActive(t.Context(), p.ID)

			require.ErrorIs(t, err, errDB)
		})
	}
}
