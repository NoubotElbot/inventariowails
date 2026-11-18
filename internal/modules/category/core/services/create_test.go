package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
)

func TestCategoryService_Create(t *testing.T) {
	t.Run("crea categoría válida", func(t *testing.T) {
		svc, repo := newTestService(t)

		got, err := svc.Create(t.Context(), ports.CreateInput{Name: "Electrónica"})

		require.NoError(t, err)
		assert.NotZero(t, got.ID)
		assert.Equal(t, "Electrónica", got.Name)
		assert.Contains(t, repo.categories, got.ID)
	})

	tests := []struct {
		name    string
		input   ports.CreateInput
		saveErr error
		wantErr error
	}{
		{name: "falla con nombre vacío", input: ports.CreateInput{Name: ""}, wantErr: domain.ErrNameRequired},
		{name: "propaga error del repositorio", input: ports.CreateInput{Name: "Electrónica"}, saveErr: errDB, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			repo.saveErr = tt.saveErr

			_, err := svc.Create(t.Context(), tt.input)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
