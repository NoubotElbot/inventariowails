package handler

import "inventario/internal/platform/pagination"

type SaleItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CreateSaleRequest struct {
	Items []SaleItemRequest `json:"items"`
}

type ListSalesRequest struct {
	pagination.Params
}
