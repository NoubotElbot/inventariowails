package category

import "inventario/internal/domain/pagination"

// FindAllQuery carries the search and pagination filters for FindAll.
type FindAllQuery struct {
	Search string
	pagination.Params
}
