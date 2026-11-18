package handler

import "inventario/internal/platform/pagination"

type ListProductsRequest struct {
	Search      string  `json:"search"`
	CategoryIDs []int64 `json:"category_ids"`
	pagination.Params
}

type CreateProductRequest struct {
	CategoryID *int64 `json:"category_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Barcode    string `json:"barcode"`
	Price      int64  `json:"price"`
	Stock      int    `json:"stock"`
	MinStock   int    `json:"min_stock"`
}

type UpdateProductRequest struct {
	CategoryID *int64 `json:"category_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Barcode    string `json:"barcode"`
	Price      int64  `json:"price"`
	MinStock   int    `json:"min_stock"`
}
