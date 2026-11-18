package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
)

func TestCategoryService_Update(t *testing.T) {
	t.Run("actualiza categoría existente", func(t *testing.T) {
		svc, repo := newTestService(t)
		c := seedCategory(t, repo, "Electrónica")

		got, err := svc.Update(t.Context(), ports.UpdateInput{ID: c.ID, Name: "Electrónica y más"})

		require.NoError(t, err)
		assert.Equal(t, "Electrónica y más", got.Name)
		assert.Equal(t, "Electrónica y más", repo.categories[c.ID].Name)
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mockRepo) int64
		newName string
		wantErr error
	}{
		{
			name:    "falla con nombre vacío",
			setup:   func(t *testing.T, r *mockRepo) int64 { return seedCategory(t, r, "Electrónica").ID },
			newName: "",
			wantErr: domain.ErrNameRequired,
		},
		{
			name:    "falla con ID inexistente",
			setup:   func(*testing.T, *mockRepo) int64 { return 999 },
			newName: "Algo",
			wantErr: domain.ErrNotFound,
		},
		{
			name: "propaga error del repositorio al buscar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.findByIDErr = errDB
				return seedCategory(t, r, "Electrónica").ID
			},
			newName: "Nuevo nombre",
			wantErr: errDB,
		},
		{
			name: "propaga error del repositorio al actualizar",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.updateErr = errDB
				return seedCategory(t, r, "Electrónica").ID
			},
			newName: "Nuevo nombre",
			wantErr: errDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			id := tt.setup(t, repo)

			_, err := svc.Update(t.Context(), ports.UpdateInput{ID: id, Name: tt.newName})

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
