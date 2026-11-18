package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/platform/pagination"
)

// seedSales guarda n ventas de un item cada una en el repositorio falso.
func seedSales(t *testing.T, d testDeps, n int) {
	t.Helper()
	for range n {
		require.NoError(t, d.sales.Save(t.Context(), &domain.Sale{
			Total: 100,
			Items: []domain.SaleItem{{ProductID: 1, Quantity: 1, UnitPrice: 100, Subtotal: 100}},
		}))
	}
}

func query(page, pageSize int) domain.FindAllQuery {
	return domain.FindAllQuery{Params: pagination.Params{Page: page, PageSize: pageSize}}
}

func TestSaleService_List(t *testing.T) {
	t.Run("retorna página vacía sin error", func(t *testing.T) {
		svc, _ := newTestService(t)

		got, err := svc.List(t.Context(), query(1, 10))

		require.NoError(t, err)
		assert.Empty(t, got.Data)
		assert.Zero(t, got.Meta.TotalItems)
	})

	t.Run("retorna las ventas del repositorio", func(t *testing.T) {
		svc, d := newTestService(t)
		seedSales(t, d, 2)

		got, err := svc.List(t.Context(), query(1, 10))

		require.NoError(t, err)
		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Meta.TotalItems)
	})

	t.Run("pagina según los parámetros", func(t *testing.T) {
		svc, d := newTestService(t)
		seedSales(t, d, 5)

		got, err := svc.List(t.Context(), query(2, 2))

		require.NoError(t, err)
		assert.Len(t, got.Data, 2)
		assert.Equal(t, 2, got.Meta.Page)
		assert.Equal(t, int64(5), got.Meta.TotalItems)
		assert.Equal(t, 3, got.Meta.TotalPages)
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		svc, d := newTestService(t)
		d.sales.findAllErr = errDB

		_, err := svc.List(t.Context(), query(1, 10))

		require.ErrorIs(t, err, errDB)
	})
}
