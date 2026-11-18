package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/product/core/domain"
)

func TestProductService_Get(t *testing.T) {
	t.Run("encuentra producto existente", func(t *testing.T) {
		svc, repo := newTestService(t)
		p := seedProduct(t, repo, "LAP-001")

		got, err := svc.Get(t.Context(), p.ID)

		require.NoError(t, err)
		assert.Equal(t, p, *got)
	})

	tests := []struct {
		name        string
		findByIDErr error
		wantErr     error
	}{
		{name: "ID inexistente retorna ErrNotFound", wantErr: domain.ErrNotFound},
		{name: "propaga error del repositorio", findByIDErr: errDB, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			repo.findByIDErr = tt.findByIDErr

			_, err := svc.Get(t.Context(), 999)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
