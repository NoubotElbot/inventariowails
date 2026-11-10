package domain_test

import (
	"inventario/internal/modules/category/domain"
	"testing"
)

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.Category
		wantErr error
	}{
		{
			name:    "válido",
			input:   domain.Category{Name: "Electrónica"},
			wantErr: nil,
		},
		{
			name:    "nombre vacío",
			input:   domain.Category{Name: ""},
			wantErr: domain.ErrNameRequired,
		},
		{
			name:    "nombre solo espacios",
			input:   domain.Category{Name: "   "},
			wantErr: domain.ErrNameRequired,
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
