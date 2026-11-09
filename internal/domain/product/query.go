package product

import "inventario/internal/domain/pagination"

// FindAllQuery carries the search, category, and pagination filters for FindAll.
type FindAllQuery struct {
	Search      string
	CategoryIDs []int64
	pagination.Params
}
