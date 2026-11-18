package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/sale/core/domain"
)

func TestSaleService_Get(t *testing.T) {
	t.Run("encuentra venta existente con items", func(t *testing.T) {
		svc, d := newTestService(t)
		s := &domain.Sale{
			Total: 1998,
			Items: []domain.SaleItem{{ProductID: 1, Quantity: 2, UnitPrice: 999, Subtotal: 1998}},
		}
		require.NoError(t, d.sales.Save(t.Context(), s))

		got, err := svc.Get(t.Context(), s.ID)

		require.NoError(t, err)
		assert.Equal(t, *s, *got)
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
			svc, d := newTestService(t)
			d.sales.findByIDErr = tt.findByIDErr

			_, err := svc.Get(t.Context(), 999)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
