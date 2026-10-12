package product_test

import (
	"inventario/internal/domain/product"
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	validProduct := func() product.Product {
		return product.Product{
			Name:  "Laptop",
			Code:  "LAP-001",
			Price: 999.99,
			Stock: 10,
		}
	}

	tests := []struct {
		name    string
		input   product.Product
		wantErr error
	}{
		{
			name:    "válido",
			input:   validProduct(),
			wantErr: nil,
		},
		{
			name: "nombre vacío",
			input: func() product.Product {
				p := validProduct()
				p.Name = ""
				return p
			}(),
			wantErr: product.ErrNameRequired,
		},
		{
			name: "nombre solo espacios",
			input: func() product.Product {
				p := validProduct()
				p.Name = "   "
				return p
			}(),
			wantErr: product.ErrNameRequired,
		},
		{
			name: "código vacío",
			input: func() product.Product {
				p := validProduct()
				p.Code = ""
				return p
			}(),
			wantErr: product.ErrCodeRequired,
		},
		{
			name: "precio cero",
			input: func() product.Product {
				p := validProduct()
				p.Price = 0
				return p
			}(),
			wantErr: product.ErrPriceInvalid,
		},
		{
			name: "precio negativo",
			input: func() product.Product {
				p := validProduct()
				p.Price = -1
				return p
			}(),
			wantErr: product.ErrPriceInvalid,
		},
		{
			name: "stock negativo",
			input: func() product.Product {
				p := validProduct()
				p.Stock = -1
				return p
			}(),
			wantErr: product.ErrStockNegative,
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

func TestProduct_HasLowStock(t *testing.T) {
	tests := []struct {
		name     string
		stock    int
		minStock int
		want     bool
	}{
		{name: "stock sobre mínimo", stock: 10, minStock: 5, want: false},
		{name: "stock igual al mínimo", stock: 5, minStock: 5, want: true},
		{name: "stock bajo el mínimo", stock: 2, minStock: 5, want: true},
		{name: "stock cero", stock: 0, minStock: 0, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := product.Product{Stock: tt.stock, MinStock: tt.minStock}
			if got := p.HasLowStock(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
