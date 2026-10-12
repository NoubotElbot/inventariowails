package stock_test

import (
	"inventario/internal/domain/stock"
	"testing"
)

func TestMovement_Validate(t *testing.T) {
	validMovement := func() stock.Movement {
		return stock.Movement{
			ProductID:    1,
			Quantity:     10,
			MovementType: stock.MovementAdjustment,
		}
	}

	tests := []struct {
		name    string
		input   stock.Movement
		wantErr error
	}{
		{
			name:    "válido adjustment",
			input:   validMovement(),
			wantErr: nil,
		},
		{
			name: "válido sale",
			input: stock.Movement{
				ProductID:    1,
				Quantity:     -5,
				MovementType: stock.MovementSale,
			},
			wantErr: nil,
		},
		{
			name: "válido return",
			input: stock.Movement{
				ProductID:    1,
				Quantity:     3,
				MovementType: stock.MovementReturn,
			},
			wantErr: nil,
		},
		{
			name: "producto requerido",
			input: func() stock.Movement {
				m := validMovement()
				m.ProductID = 0
				return m
			}(),
			wantErr: stock.ErrProductRequired,
		},
		{
			name: "cantidad cero",
			input: func() stock.Movement {
				m := validMovement()
				m.Quantity = 0
				return m
			}(),
			wantErr: stock.ErrQuantityZero,
		},
		{
			name: "tipo de movimiento inválido",
			input: func() stock.Movement {
				m := validMovement()
				m.MovementType = "unknown"
				return m
			}(),
			wantErr: stock.ErrInvalidMovement,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
