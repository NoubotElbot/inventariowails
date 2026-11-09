package domain

import "inventario/internal/platform/pagination"

// FindAllQuery carries the search, category, and pagination filters for FindAll.
type FindAllQuery struct {
	Search      string
	CategoryIDs []int64
	pagination.Params
}
