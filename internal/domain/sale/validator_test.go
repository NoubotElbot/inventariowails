package sale_test

import (
	"inventario/internal/domain/sale"
	"testing"
)

func TestSale_Validate(t *testing.T) {
	validSale := func() sale.Sale {
		return sale.Sale{
			Items: []sale.SaleItem{
				{ProductID: 1, Quantity: 2, UnitPrice: 100, Subtotal: 200},
			},
		}
	}

	tests := []struct {
		name    string
		input   sale.Sale
		wantErr error
	}{
		{
			name:    "válida con un item",
			input:   validSale(),
			wantErr: nil,
		},
		{
			name: "válida con múltiples items",
			input: sale.Sale{
				Items: []sale.SaleItem{
					{ProductID: 1, Quantity: 2, UnitPrice: 100, Subtotal: 200},
					{ProductID: 2, Quantity: 1, UnitPrice: 50, Subtotal: 50},
				},
			},
			wantErr: nil,
		},
		{
			name:    "sin items",
			input:   sale.Sale{},
			wantErr: sale.ErrNoItems,
		},
		{
			name: "item sin producto",
			input: sale.Sale{
				Items: []sale.SaleItem{
					{ProductID: 0, Quantity: 1},
				},
			},
			wantErr: sale.ErrProductInvalid,
		},
		{
			name: "item con cantidad cero",
			input: sale.Sale{
				Items: []sale.SaleItem{
					{ProductID: 1, Quantity: 0},
				},
			},
			wantErr: sale.ErrQuantityInvalid,
		},
		{
			name: "item con cantidad negativa",
			input: sale.Sale{
				Items: []sale.SaleItem{
					{ProductID: 1, Quantity: -1},
				},
			},
			wantErr: sale.ErrQuantityInvalid,
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
