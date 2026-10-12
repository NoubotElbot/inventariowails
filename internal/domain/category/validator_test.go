package category_test

import (
	"inventario/internal/domain/category"
	"testing"
)

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   category.Category
		wantErr error
	}{
		{
			name:    "válido",
			input:   category.Category{Name: "Electrónica"},
			wantErr: nil,
		},
		{
			name:    "nombre vacío",
			input:   category.Category{Name: ""},
			wantErr: category.ErrNameRequired,
		},
		{
			name:    "nombre solo espacios",
			input:   category.Category{Name: "   "},
			wantErr: category.ErrNameRequired,
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
