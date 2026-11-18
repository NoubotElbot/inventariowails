package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
)

func TestStockService_AdjustStock(t *testing.T) {
	valid := []struct {
		name      string
		input     ports.AdjustInput
		wantStock int
	}{
		{"ajuste de entrada", ports.AdjustInput{Quantity: 3, Type: domain.MovementAdjustmentIn}, 13},
		{"ajuste de salida con nota", ports.AdjustInput{Quantity: 1, Type: domain.MovementAdjustmentOut, Notes: "conteo"}, 9},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			id := d.products.seed(10)
			tt.input.ProductID = id

			got, err := svc.AdjustStock(t.Context(), tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStock, d.products.stockOf(id))
			assert.Equal(t, tt.input.Type, got.MovementType)
			assert.Equal(t, tt.input.Notes, got.Notes)
			require.Len(t, d.stock.movements, 1)
		})
	}

	// Solo se aceptan ajustes manuales; las validaciones del dominio se
	// evalúan primero y en orden (producto, cantidad, tipo).
	invalid := []struct {
		name    string
		input   ports.AdjustInput
		wantErr error
	}{
		{"sin producto", ports.AdjustInput{Quantity: 1, Type: domain.MovementAdjustmentIn}, domain.ErrProductRequired},
		{"cantidad cero", ports.AdjustInput{ProductID: 1, Type: domain.MovementAdjustmentIn}, domain.ErrQuantityNotPositive},
		{"cantidad negativa", ports.AdjustInput{ProductID: 1, Quantity: -2, Type: domain.MovementAdjustmentIn}, domain.ErrQuantityNotPositive},
		{"venta rechazada", ports.AdjustInput{ProductID: 1, Quantity: 1, Type: domain.MovementSale}, domain.ErrInvalidMovement},
		{"devolución rechazada", ports.AdjustInput{ProductID: 1, Quantity: 1, Type: domain.MovementReturn}, domain.ErrInvalidMovement},
		{"tipo vacío", ports.AdjustInput{ProductID: 1, Quantity: 1}, domain.ErrInvalidMovement},
		{"tipo desconocido", ports.AdjustInput{ProductID: 1, Quantity: 1, Type: "foo"}, domain.ErrInvalidMovement},
		{"varios fallos devuelve el primero", ports.AdjustInput{}, domain.ErrProductRequired},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			d.products.seed(10)

			_, err := svc.AdjustStock(t.Context(), tt.input)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Zero(t, d.uow.calls, "no debe abrir transacción")
			assert.Empty(t, d.stock.movements)
			assert.Equal(t, 10, d.products.stockOf(1))
		})
	}

	t.Run("propaga errores del registro del movimiento", func(t *testing.T) {
		svc, _ := newTestService(t)

		_, err := svc.AdjustStock(t.Context(), ports.AdjustInput{
			ProductID: 999, Quantity: 1, Type: domain.MovementAdjustmentIn,
		})

		require.Error(t, err)
	})
}
