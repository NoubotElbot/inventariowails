package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/platform/pagination"
)

// El filtrado y la paginación los hace el repositorio; aquí solo se verifica
// que el servicio pase la query intacta y devuelva o envuelva lo que recibe.
func TestCategoryService_List(t *testing.T) {
	t.Run("reenvía la query sin modificarla", func(t *testing.T) {
		queries := map[string]domain.FindAllQuery{
			"vacía":        {},
			"paginada":     {Params: pagination.Params{Page: 3, PageSize: 2}},
			"con búsqueda": {Search: "elec"},
			"solo inactivas": {
				Params:       pagination.Params{Page: 1, PageSize: 10},
				OnlyInactive: true,
			},
		}
		for name, q := range queries {
			t.Run(name, func(t *testing.T) {
				svc, repo := newTestService(t)

				_, err := svc.List(t.Context(), q)

				require.NoError(t, err)
				assert.Equal(t, q, repo.lastQuery)
			})
		}
	})

	t.Run("retorna el resultado del repositorio", func(t *testing.T) {
		svc, repo := newTestService(t)
		want := pagination.Result[domain.Category]{
			Data: []domain.Category{{ID: 1, Name: "Electrónica"}, {ID: 2, Name: "Ropa"}},
			Meta: pagination.Meta{
				Page:       pagination.DefaultPage,
				PageSize:   pagination.DefaultPageSize,
				TotalItems: 2,
				TotalPages: 1,
			},
		}
		repo.findAllResult = want

		got, err := svc.List(t.Context(), domain.FindAllQuery{})

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		svc, repo := newTestService(t)
		repo.findAllErr = errDB

		_, err := svc.List(t.Context(), domain.FindAllQuery{})

		require.ErrorIs(t, err, errDB)
	})
}
