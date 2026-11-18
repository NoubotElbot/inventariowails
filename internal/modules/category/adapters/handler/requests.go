package handler

import (
	"inventario/internal/platform/pagination"
)

// ListCategoriesRequest carries pagination and optional search filter.
type ListCategoriesRequest struct {
	Search       string `json:"search"`
	OnlyInactive bool   `json:"only_inactive"`
	pagination.Params
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}
