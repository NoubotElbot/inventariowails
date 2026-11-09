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
	items := make([]R, 0, len(page.Items))
	for _, it := range page.Items {
		items = append(items, fn(it))
	}
	return PageResponse[R]{
		Items:      items,
		Total:      page.Total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
	}
}
