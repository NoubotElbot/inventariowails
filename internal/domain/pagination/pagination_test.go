package pagination_test

import (
	"testing"

	"inventario/internal/domain/pagination"
)

func TestParams_Normalize(t *testing.T) {
	tests := []struct {
		name string
		in   pagination.Params
		want pagination.Params
	}{
		{"zero values", pagination.Params{}, pagination.Params{Page: 1, PageSize: 20}},
		{"negative", pagination.Params{Page: -3, PageSize: -1}, pagination.Params{Page: 1, PageSize: 20}},
		{"valid", pagination.Params{Page: 3, PageSize: 10}, pagination.Params{Page: 3, PageSize: 10}},
		{"oversize", pagination.Params{Page: 1, PageSize: 1000}, pagination.Params{Page: 1, PageSize: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Normalize(); got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParams_OffsetLimit(t *testing.T) {
	p := pagination.Params{Page: 3, PageSize: 10}
	if p.Offset() != 20 || p.Limit() != 10 {
		t.Errorf("offset=%d limit=%d, want 20/10", p.Offset(), p.Limit())
	}
	if (pagination.Params{}).Offset() != 0 {
		t.Error("zero params should have offset 0")
	}
}

func TestNewResult(t *testing.T) {
	tests := []struct {
		name  string
		total int64
		size  int
		want  int
	}{
		{"empty", 0, 10, 0},
		{"exact multiple", 20, 10, 2},
		{"remainder", 21, 10, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pagination.NewResult[int](nil, tt.total, pagination.Params{Page: 1, PageSize: tt.size})
			if r.TotalPages != tt.want {
				t.Errorf("TotalPages=%d, want %d", r.TotalPages, tt.want)
			}
			if r.Items == nil {
				t.Error("Items must be non-nil")
			}
		})
	}
}
