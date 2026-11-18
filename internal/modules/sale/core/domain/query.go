package domain

import "inventario/internal/platform/pagination"

// FindAllQuery carries the pagination filters for FindAll.
type FindAllQuery struct {
	pagination.Params
}
