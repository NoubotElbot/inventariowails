package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/stock/core/domain"
)

func TestStockService_ListMovements(t *testing.T) {
	t.Run("retorna solo los movimientos del producto", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(10)
		other := d.products.seed(10)
		require.NoError(t, d.stock.Save(t.Context(), &domain.Movement{ProductID: id, Quantity: 5, MovementType: domain.MovementAdjustmentIn}))
		require.NoError(t, d.stock.Save(t.Context(), &domain.Movement{ProductID: id, Quantity: 2, MovementType: domain.MovementSale}))
		require.NoError(t, d.stock.Save(t.Context(), &domain.Movement{ProductID: other, Quantity: 1, MovementType: domain.MovementReturn}))

		got, err := svc.ListMovements(t.Context(), id)

		require.NoError(t, err)
		require.Len(t, got, 2)
		for _, mv := range got {
			assert.Equal(t, id, mv.ProductID)
		}
	})

	t.Run("retorna lista vacía sin movimientos", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(10)

		got, err := svc.ListMovements(t.Context(), id)

		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})

	tests := []struct {
		name      string
		productID int64
		setup     func(d testDeps)
		wantErr   error
	}{
		{name: "producto cero devuelve ErrProductRequired", productID: 0, wantErr: domain.ErrProductRequired},
		{name: "producto inexistente devuelve ErrNotFound", productID: 999, wantErr: productdomain.ErrNotFound},
		{name: "propaga error al buscar producto", productID: 1, setup: func(d testDeps) { d.products.findByIDErr = errDB }, wantErr: errDB},
		{name: "propaga error del repositorio de movimientos", productID: 1, setup: func(d testDeps) { d.stock.findErr = errDB }, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			d.products.seed(10)
			if tt.setup != nil {
				tt.setup(d)
			}

			_, err := svc.ListMovements(t.Context(), tt.productID)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
