package handler

import (
	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/platform/responses"
)

type SaleProductResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type SaleItemResponse struct {
	ID        int64                `json:"id"`
	ProductID int64                `json:"product_id"`
	Product   *SaleProductResponse `json:"product"`
	Quantity  int                  `json:"quantity"`
	UnitPrice int64                `json:"unit_price"`
	Subtotal  int64                `json:"subtotal"`
	CreatedAt string               `json:"created_at"`
}

type SaleResponse struct {
	ID        int64              `json:"id"`
	Total     int64              `json:"total"`
	Items     []SaleItemResponse `json:"items"`
	CreatedAt string             `json:"created_at"`
}

func toResponse(s domain.Sale) SaleResponse {
	return SaleResponse{
		ID:        s.ID,
		Total:     s.Total,
		Items:     responses.MapSlice(s.Items, toItemResponse),
		CreatedAt: responses.FormatTime(s.CreatedAt),
	}
}

func toItemResponse(item domain.SaleItem) SaleItemResponse {
	return SaleItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Product:   toProductResponse(item.Product),
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Subtotal:  item.Subtotal,
		CreatedAt: responses.FormatTime(item.CreatedAt),
	}
}

func toProductResponse(p *domain.SaleProduct) *SaleProductResponse {
	if p == nil {
		return nil
	}
	return &SaleProductResponse{ID: p.ID, Name: p.Name, Code: p.Code}
}
