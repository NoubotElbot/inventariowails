package handler

import (
	"inventario/internal/modules/product/core/domain"
	"inventario/internal/platform/responses"
)

type ProductResponse struct {
	ID         int64  `json:"id"`
	CategoryID *int64 `json:"category_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Barcode    string `json:"barcode"`
	Price      int64  `json:"price"`
	Stock      int    `json:"stock"`
	MinStock   int    `json:"min_stock"`
	Active     bool   `json:"active"`
	LowStock   bool   `json:"low_stock"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func toResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		CategoryID: p.CategoryID,
		Name:       p.Name,
		Code:       p.Code,
		Barcode:    p.Barcode,
		Price:      p.Price,
		Stock:      p.Stock,
		MinStock:   p.MinStock,
		Active:     p.Active,
		LowStock:   p.HasLowStock(),
		CreatedAt:  responses.FormatTime(p.CreatedAt),
		UpdatedAt:  responses.FormatTime(p.UpdatedAt),
	}
}
