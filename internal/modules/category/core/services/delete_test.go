package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
)

func TestCategoryService_Delete(t *testing.T) {
	t.Run("elimina suavemente categoría existente", func(t *testing.T) {
		svc, repo := newTestService(t)
		c := seedCategory(t, repo, "Electrónica")

		require.NoError(t, svc.Delete(t.Context(), c.ID))

		assert.NotNil(t, repo.categories[c.ID].DeletedAt)
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
			name:    "falla si la categoría ya está eliminada",
			setup:   func(t *testing.T, r *mockRepo) int64 { return seedDeletedCategory(t, r, "Electrónica").ID },
			wantErr: domain.ErrAlreadyDeleted,
		},
		{
			name: "propaga error del repositorio al buscar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.findByIDErr = errDB
				return seedCategory(t, r, "Electrónica").ID
			},
			wantErr: errDB,
		},
		{
			name: "propaga error del repositorio al eliminar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.softDeleteErr = errDB
				return seedCategory(t, r, "Electrónica").ID
			},
			wantErr: errDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			id := tt.setup(t, repo)

			require.ErrorIs(t, svc.Delete(t.Context(), id), tt.wantErr)
		})
	}
}
