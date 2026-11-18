package responses

import "inventario/internal/platform/pagination"

// PageResponse is the shared paginated DTO returned to the frontend.
type PageResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// MapPage converts a domain page into a PageResponse using fn for each item.
func MapPage[D any, R any](page pagination.Result[D], fn func(D) R) PageResponse[R] {
	items := MapSlice(page.Data, fn)
	return PageResponse[R]{
		Items:      items,
		Total:      page.Meta.TotalItems,
		Page:       page.Meta.Page,
		PageSize:   page.Meta.PageSize,
		TotalPages: page.Meta.TotalPages,
	}
}
