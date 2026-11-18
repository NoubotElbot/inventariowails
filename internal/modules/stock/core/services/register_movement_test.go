package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
)

func TestStockService_RegisterMovement(t *testing.T) {
	// Cada tipo con cantidad positiva mueve el stock en su dirección y se
	// guarda con la cantidad positiva.
	directions := []struct {
		typ       domain.MovementType
		wantStock int
	}{
		{domain.MovementSale, 7},
		{domain.MovementAdjustmentOut, 7},
		{domain.MovementReturn, 13},
		{domain.MovementAdjustmentIn, 13},
	}
	for _, tt := range directions {
		t.Run("cantidad positiva con tipo "+string(tt.typ), func(t *testing.T) {
			svc, d := newTestService(t)
			id := d.products.seed(10)

			got, err := svc.RegisterMovement(t.Context(), ports.RegisterMovementInput{
				ProductID: id, Quantity: 3, Type: tt.typ, Notes: "nota",
			})

			require.NoError(t, err)
			assert.NotZero(t, got.ID)
			assert.False(t, got.CreatedAt.IsZero())
			assert.Equal(t, tt.wantStock, d.products.stockOf(id))
			require.Len(t, d.stock.movements, 1)
			mv := d.stock.movements[0]
			assert.Equal(t, tt.typ, mv.MovementType)
			assert.Equal(t, 3, mv.Quantity)
			assert.Equal(t, "nota", mv.Notes)
			assert.Equal(t, 1, d.uow.calls, "una transacción por movimiento")
		})
	}

	t.Run("stock insuficiente no registra movimiento ni cambia stock", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(2)

		_, err := svc.RegisterMovement(t.Context(), ports.RegisterMovementInput{
			ProductID: id, Quantity: 3, Type: domain.MovementSale,
		})

		require.ErrorIs(t, err, productdomain.ErrInsufficientStock)
		assert.Empty(t, d.stock.movements)
		assert.Equal(t, 2, d.products.stockOf(id))
	})

	t.Run("permite dejar el stock en cero", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(3)

		_, err := svc.RegisterMovement(t.Context(), ports.RegisterMovementInput{
			ProductID: id, Quantity: 3, Type: domain.MovementSale,
		})

		require.NoError(t, err)
		assert.Zero(t, d.products.stockOf(id))
	})

	t.Run("falla con producto inexistente", func(t *testing.T) {
		svc, _ := newTestService(t)

		_, err := svc.RegisterMovement(t.Context(), ports.RegisterMovementInput{
			ProductID: 999, Quantity: 1, Type: domain.MovementSale,
		})

		require.ErrorIs(t, err, productdomain.ErrNotFound)
	})

	invalid := []struct {
		name    string
		input   ports.RegisterMovementInput
		wantErr error
	}{
		{"sin producto", ports.RegisterMovementInput{Quantity: 1, Type: domain.MovementSale}, domain.ErrProductRequired},
		{"cantidad cero", ports.RegisterMovementInput{ProductID: 1, Quantity: 0, Type: domain.MovementSale}, domain.ErrQuantityNotPositive},
		{"cantidad negativa", ports.RegisterMovementInput{ProductID: 1, Quantity: -3, Type: domain.MovementSale}, domain.ErrQuantityNotPositive},
		{"tipo desconocido", ports.RegisterMovementInput{ProductID: 1, Quantity: 1, Type: "robo"}, domain.ErrInvalidMovement},
		{"tipo adjustment antiguo", ports.RegisterMovementInput{ProductID: 1, Quantity: 1, Type: "adjustment"}, domain.ErrInvalidMovement},
	}
	for _, tt := range invalid {
		t.Run(tt.name+" falla sin abrir transacción", func(t *testing.T) {
			svc, d := newTestService(t)
			d.products.seed(10)

			_, err := svc.RegisterMovement(t.Context(), tt.input)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Zero(t, d.uow.calls)
			assert.Empty(t, d.stock.movements)
		})
	}

	repoErrors := []struct {
		name  string
		setup func(d testDeps)
	}{
		{"propaga error al buscar producto", func(d testDeps) { d.products.findByIDErr = errDB }},
		{"propaga error al actualizar stock", func(d testDeps) { d.products.updateStockErr = errDB }},
		{"propaga error al guardar movimiento", func(d testDeps) { d.stock.saveErr = errDB }},
		{"propaga error de la unidad de trabajo", func(d testDeps) { d.uow.err = errDB }},
	}
	for _, tt := range repoErrors {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			id := d.products.seed(10)
			tt.setup(d)

			_, err := svc.RegisterMovement(t.Context(), ports.RegisterMovementInput{
				ProductID: id, Quantity: 1, Type: domain.MovementSale,
			})

			require.ErrorIs(t, err, errDB)
		})
	}
}
