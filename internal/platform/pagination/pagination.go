// Package pagination provides shared, domain-agnostic pagination types.
package pagination

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Params are the pagination inputs coming from the caller.
type Params struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Normalize returns params with safe values applied.
func (p Params) Normalize() Params {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p
}

// Limit returns the SQL LIMIT for the (normalized) params.
func (p Params) Limit() int { return p.Normalize().PageSize }

// Offset returns the SQL OFFSET for the (normalized) params.
func (p Params) Offset() int {
	n := p.Normalize()
	return (n.Page - 1) * n.PageSize
}

// Result is a page of items plus the metadata needed to navigate.
type Result[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewResult builds a Result, normalizing params and guaranteeing non-nil Items.
func NewResult[T any](items []T, total int64, p Params) Result[T] {
	p = p.Normalize()
	if items == nil {
		items = []T{}
	}
	totalPages := int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	return Result[T]{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	}
}
