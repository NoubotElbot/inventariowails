package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	stockdomain "inventario/internal/modules/stock/core/domain"
)

func TestSaleService_Create(t *testing.T) {
	t.Run("crea venta con un item", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(999, 10, true)

		got, err := svc.Create(t.Context(), items(ports.SaleItemInput{ProductID: id, Quantity: 2}))

		require.NoError(t, err)
		assert.NotZero(t, got.ID)
		assert.False(t, got.CreatedAt.IsZero())
		require.Len(t, got.Items, 1)
		assert.Equal(t, domain.SaleItem{
			ID: 1, SaleID: got.ID, ProductID: id, Quantity: 2, UnitPrice: 999, Subtotal: 1998,
			CreatedAt: got.Items[0].CreatedAt,
		}, got.Items[0])
		assert.Equal(t, int64(1998), got.Total)
		assert.Contains(t, d.sales.sales, got.ID)
	})

	t.Run("crea venta con múltiples items y suma el total", func(t *testing.T) {
		svc, d := newTestService(t)
		laptop := d.products.seed(999, 10, true)
		mouse := d.products.seed(29, 5, true)

		got, err := svc.Create(t.Context(), items(
			ports.SaleItemInput{ProductID: laptop, Quantity: 1},
			ports.SaleItemInput{ProductID: mouse, Quantity: 2},
		))

		require.NoError(t, err)
		require.Len(t, got.Items, 2)
		assert.Equal(t, int64(58), got.Items[1].Subtotal)
		assert.Equal(t, int64(999+29*2), got.Total)
	})

	t.Run("descuenta stock de cada item como movimiento sale", func(t *testing.T) {
		svc, d := newTestService(t)
		id1 := d.products.seed(10, 10, true)
		id2 := d.products.seed(10, 10, true)

		got, err := svc.Create(t.Context(), items(
			ports.SaleItemInput{ProductID: id1, Quantity: 3},
			ports.SaleItemInput{ProductID: id2, Quantity: 1},
		))

		require.NoError(t, err)
		require.Len(t, d.stock.calls, 2)
		notes := fmt.Sprintf("venta #%d", got.ID)
		for i, want := range []struct {
			id  int64
			qty int
		}{{id1, 3}, {id2, 1}} {
			mv := d.stock.calls[i]
			assert.Equal(t, want.id, mv.ProductID)
			assert.Equal(t, want.qty, mv.Quantity, "cantidad positiva; el tipo sale indica salida")
			assert.Equal(t, stockdomain.MovementSale, mv.Type)
			assert.Equal(t, notes, mv.Notes)
		}
	})

	t.Run("ejecuta todo dentro de una sola transacción", func(t *testing.T) {
		svc, d := newTestService(t)
		id := d.products.seed(10, 10, true)

		_, err := svc.Create(t.Context(), items(ports.SaleItemInput{ProductID: id, Quantity: 2}))

		require.NoError(t, err)
		assert.Equal(t, 1, d.uow.calls)
	})

	// Fallos detectados antes de la transacción: no se abre, no se guarda la
	// venta ni se mueve stock.
	preChecks := []struct {
		name    string
		input   func(d testDeps) ports.CreateInput
		setup   func(d testDeps)
		wantErr error
	}{
		{
			name:    "sin items",
			input:   func(testDeps) ports.CreateInput { return ports.CreateInput{} },
			wantErr: domain.ErrNoItems,
		},
		{
			name:    "producto inexistente",
			input:   func(testDeps) ports.CreateInput { return items(ports.SaleItemInput{ProductID: 999, Quantity: 1}) },
			wantErr: productdomain.ErrNotFound,
		},
		{
			name: "producto inactivo",
			input: func(d testDeps) ports.CreateInput {
				return items(ports.SaleItemInput{ProductID: d.products.seed(10, 10, false), Quantity: 1})
			},
			wantErr: domain.ErrProductInactive,
		},
		{
			name: "stock insuficiente",
			input: func(d testDeps) ports.CreateInput {
				return items(ports.SaleItemInput{ProductID: d.products.seed(10, 2, true), Quantity: 5})
			},
			wantErr: productdomain.ErrInsufficientStock,
		},
		{
			name: "cantidad cero",
			input: func(d testDeps) ports.CreateInput {
				return items(ports.SaleItemInput{ProductID: d.products.seed(10, 10, true), Quantity: 0})
			},
			wantErr: domain.ErrQuantityInvalid,
		},
		{
			name: "propaga error al buscar producto",
			input: func(d testDeps) ports.CreateInput {
				return items(ports.SaleItemInput{ProductID: d.products.seed(10, 10, true), Quantity: 1})
			},
			setup:   func(d testDeps) { d.products.findByIDErr = errDB },
			wantErr: errDB,
		},
	}
	for _, tt := range preChecks {
		t.Run(tt.name+" falla sin abrir transacción", func(t *testing.T) {
			svc, d := newTestService(t)
			input := tt.input(d)
			if tt.setup != nil {
				tt.setup(d)
			}

			got, err := svc.Create(t.Context(), input)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, got)
			assert.Zero(t, d.uow.calls)
			assert.Empty(t, d.sales.sales)
			assert.Empty(t, d.stock.calls)
		})
	}

	// Fallos dentro de la transacción: se propagan sin envolver la causa.
	inTx := []struct {
		name    string
		setup   func(d testDeps)
		wantErr error
	}{
		{"propaga error al guardar la venta", func(d testDeps) { d.sales.saveErr = errDB }, errDB},
		{"propaga error del registro de stock", func(d testDeps) { d.stock.err = errDB }, errDB},
		{
			// El stock cambió entre la verificación previa y la transacción.
			"propaga stock insuficiente detectado por el módulo de stock",
			func(d testDeps) { d.stock.err = productdomain.ErrInsufficientStock },
			productdomain.ErrInsufficientStock,
		},
		{"propaga error de la unidad de trabajo", func(d testDeps) { d.uow.err = errDB }, errDB},
	}
	for _, tt := range inTx {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			id := d.products.seed(10, 10, true)
			tt.setup(d)

			got, err := svc.Create(t.Context(), items(ports.SaleItemInput{ProductID: id, Quantity: 1}))

			require.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, got)
		})
	}
}
