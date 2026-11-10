package domain_test

import (
	"inventario/internal/modules/stock/domain"
	"testing"
)

func TestMovement_Validate(t *testing.T) {
	validMovement := func() domain.Movement {
		return domain.Movement{
			ProductID:    1,
			Quantity:     10,
			MovementType: domain.MovementAdjustment,
		}
	}

	tests := []struct {
		name    string
		input   domain.Movement
		wantErr error
	}{
		{
			name:    "válido adjustment",
			input:   validMovement(),
			wantErr: nil,
		},
		{
			name: "válido sale",
			input: domain.Movement{
				ProductID:    1,
				Quantity:     -5,
				MovementType: domain.MovementSale,
			},
			wantErr: nil,
		},
		{
			name: "válido return",
			input: domain.Movement{
				ProductID:    1,
				Quantity:     3,
				MovementType: domain.MovementReturn,
			},
			wantErr: nil,
		},
		{
			name: "producto requerido",
			input: func() domain.Movement {
				m := validMovement()
				m.ProductID = 0
				return m
			}(),
			wantErr: domain.ErrProductRequired,
		},
		{
			name: "cantidad cero",
			input: func() domain.Movement {
				m := validMovement()
				m.Quantity = 0
				return m
			}(),
			wantErr: domain.ErrQuantityZero,
		},
		{
			name: "tipo de movimiento inválido",
			input: func() domain.Movement {
				m := validMovement()
				m.MovementType = "unknown"
				return m
			}(),
			wantErr: domain.ErrInvalidMovement,
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
