package domain_test

import (
	"inventario/internal/modules/sale/domain"
	"testing"
)

func TestSale_Validate(t *testing.T) {
	validSale := func() domain.Sale {
		return domain.Sale{
			Items: []domain.SaleItem{
				{ProductID: 1, Quantity: 2, UnitPrice: 100, Subtotal: 200},
			},
		}
	}

	tests := []struct {
		name    string
		input   domain.Sale
		wantErr error
	}{
		{
			name:    "válida con un item",
			input:   validSale(),
			wantErr: nil,
		},
		{
			name: "válida con múltiples items",
			input: domain.Sale{
				Items: []domain.SaleItem{
					{ProductID: 1, Quantity: 2, UnitPrice: 100, Subtotal: 200},
					{ProductID: 2, Quantity: 1, UnitPrice: 50, Subtotal: 50},
				},
			},
			wantErr: nil,
		},
		{
			name:    "sin items",
			input:   domain.Sale{},
			wantErr: domain.ErrNoItems,
		},
		{
			name: "item sin producto",
			input: domain.Sale{
				Items: []domain.SaleItem{
					{ProductID: 0, Quantity: 1},
				},
			},
			wantErr: domain.ErrProductInvalid,
		},
		{
			name: "item con cantidad cero",
			input: domain.Sale{
				Items: []domain.SaleItem{
					{ProductID: 1, Quantity: 0},
				},
			},
			wantErr: domain.ErrQuantityInvalid,
		},
		{
			name: "item con cantidad negativa",
			input: domain.Sale{
				Items: []domain.SaleItem{
					{ProductID: 1, Quantity: -1},
				},
			},
			wantErr: domain.ErrQuantityInvalid,
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
