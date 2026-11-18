package handler

import (
	"inventario/internal/modules/category/core/domain"
	"inventario/internal/platform/responses"
)

type CategoryResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type CategoryOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func toResponse(c domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: responses.FormatTime(c.CreatedAt),
		UpdatedAt: responses.FormatTime(c.UpdatedAt),
		DeletedAt: responses.FormatTimePtr(c.DeletedAt),
	}
}
