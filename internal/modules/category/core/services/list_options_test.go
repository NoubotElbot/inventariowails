package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
)

func TestCategoryService_ListOptions(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		svc, _ := newTestService(t)

		got, err := svc.ListOptions(t.Context())

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("retorna las opciones del repositorio", func(t *testing.T) {
		svc, repo := newTestService(t)
		ropa := seedCategory(t, repo, "Ropa")
		elec := seedCategory(t, repo, "Electrónica")

		got, err := svc.ListOptions(t.Context())

		require.NoError(t, err)
		assert.Equal(t, []domain.Option{
			{ID: elec.ID, Name: "Electrónica"},
			{ID: ropa.ID, Name: "Ropa"},
		}, got)
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		svc, repo := newTestService(t)
		repo.findOptionsErr = errDB

		_, err := svc.ListOptions(t.Context())

		require.ErrorIs(t, err, errDB)
	})
}
