package domain

import "inventario/internal/platform/pagination"

// FindAllQuery carries the search and pagination filters for FindAll.
type FindAllQuery struct {
	Search      string
	OnlyInactive bool
	pagination.Params
}
