package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/product/core/domain"
)

func TestProductService_Delete(t *testing.T) {
	t.Run("elimina producto existente", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")

		err := svc.Delete(t.Context(), p.ID)

		require.NoError(t, err)
		assert.NotContains(t, repo.products, p.ID)
	})

	t.Run("falla con ID inexistente", func(t *testing.T) {
		svc, _ := newTestService(t)

		err := svc.Delete(t.Context(), 999)

		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	tests := []struct {
		name        string
		findByIDErr error
		deleteErr   error
	}{
		{name: "propaga error al buscar", findByIDErr: errDB},
		{name: "propaga error del repositorio en delete", deleteErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			p := seedProduct(t, repo, "LAP-001")
			repo.findByIDErr = tt.findByIDErr
			repo.deleteErr = tt.deleteErr

			err := svc.Delete(t.Context(), p.ID)

			require.ErrorIs(t, err, errDB)
		})
	}
}
