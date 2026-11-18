package domain_test

import (
	"inventario/internal/modules/stock/core/domain"
	"testing"
)

func TestMovement_Validate(t *testing.T) {
	movement := func(qty int, typ domain.MovementType) domain.Movement {
		return domain.Movement{ProductID: 1, Quantity: qty, MovementType: typ}
	}

	tests := []struct {
		name    string
		input   domain.Movement
		wantErr error
	}{
		{"válido sale", movement(5, domain.MovementSale), nil},
		{"válido return", movement(3, domain.MovementReturn), nil},
		{"válido adjustment_in", movement(2, domain.MovementAdjustmentIn), nil},
		{"válido adjustment_out", movement(2, domain.MovementAdjustmentOut), nil},
		{
			name:    "producto requerido",
			input:   domain.Movement{Quantity: 1, MovementType: domain.MovementSale},
			wantErr: domain.ErrProductRequired,
		},
		{"cantidad cero", movement(0, domain.MovementSale), domain.ErrQuantityNotPositive},
		{"cantidad negativa", movement(-5, domain.MovementSale), domain.ErrQuantityNotPositive},
		{"tipo desconocido", movement(1, "unknown"), domain.ErrInvalidMovement},
		{"tipo adjustment antiguo ya no es válido", movement(1, "adjustment"), domain.ErrInvalidMovement},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.input.Validate(); err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovement_Delta(t *testing.T) {
	tests := []struct {
		typ  domain.MovementType
		want int
	}{
		{domain.MovementSale, -4},
		{domain.MovementAdjustmentOut, -4},
		{domain.MovementReturn, 4},
		{domain.MovementAdjustmentIn, 4},
	}

	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			m := domain.Movement{ProductID: 1, Quantity: 4, MovementType: tt.typ}
			if got := m.Delta(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMovementType_IsManualAdjustment(t *testing.T) {
	tests := []struct {
		typ  domain.MovementType
		want bool
	}{
		{domain.MovementAdjustmentIn, true},
		{domain.MovementAdjustmentOut, true},
		{domain.MovementSale, false},
		{domain.MovementReturn, false},
		{"", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			if got := tt.typ.IsManualAdjustment(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
