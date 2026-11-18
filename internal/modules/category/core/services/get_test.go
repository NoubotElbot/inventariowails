package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
)

func TestCategoryService_Get(t *testing.T) {
	t.Run("encuentra categoría existente", func(t *testing.T) {
		svc, repo := newTestService(t)
		c := seedCategory(t, repo, "Electrónica")

		got, err := svc.Get(t.Context(), c.ID)

		require.NoError(t, err)
		assert.Equal(t, c, *got)
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mockRepo) int64
		wantErr error
	}{
		{
			name:    "retorna ErrNotFound para ID inexistente",
			setup:   func(*testing.T, *mockRepo) int64 { return 999 },
			wantErr: domain.ErrNotFound,
		},
		{
			name: "propaga error del repositorio",
			setup: func(t *testing.T, r *mockRepo) int64 {
				r.findByIDErr = errDB
				return seedCategory(t, r, "Electrónica").ID
			},
			wantErr: errDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			id := tt.setup(t, repo)

			_, err := svc.Get(t.Context(), id)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
