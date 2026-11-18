package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
)

func TestCategoryService_Restore(t *testing.T) {
	t.Run("restaura categoría eliminada", func(t *testing.T) {
		svc, repo := newTestService(t)
		c := seedDeletedCategory(t, repo, "Electrónica")

		require.NoError(t, svc.Restore(t.Context(), c.ID))

		assert.Nil(t, repo.categories[c.ID].DeletedAt)
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mockRepo) int64
		wantErr error
	}{
		{
			name:    "falla con ID inexistente",
			setup:   func(*testing.T, *mockRepo) int64 { return 999 },
			wantErr: domain.ErrNotFound,
		},
		{
			name:    "falla si la categoría no está eliminada",
			setup:   func(t *testing.T, r *mockRepo) int64 { return seedCategory(t, r, "Electrónica").ID },
			wantErr: domain.ErrNotDeleted,
		},
		{
			name: "propaga error del repositorio al buscar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.findByIDErr = errDB
				return seedDeletedCategory(t, r, "Electrónica").ID
			},
			wantErr: errDB,
		},
		{
			name: "propaga error del repositorio al restaurar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.restoreErr = errDB
				return seedDeletedCategory(t, r, "Electrónica").ID
			},
			wantErr: errDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			id := tt.setup(t, repo)

			require.ErrorIs(t, svc.Restore(t.Context(), id), tt.wantErr)
		})
	}
}
