// Package handler is the Wails driving adapter for the sale module.
// It translates frontend requests into calls on the ports.SaleService
// inbound port and maps domain results into response DTOs.
package handler

import (
	"context"

	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	"inventario/internal/platform/responses"
)

type SaleHandler struct {
	svc ports.SaleService
}

func NewSaleHandler(svc ports.SaleService) *SaleHandler {
	return &SaleHandler{svc: svc}
}

func (h *SaleHandler) CreateSale(ctx context.Context, req CreateSaleRequest) (*SaleResponse, error) {
	items := responses.MapSlice(req.Items, func(item SaleItemRequest) ports.SaleItemInput {
		return ports.SaleItemInput{ProductID: item.ProductID, Quantity: item.Quantity}
	})

	s, err := h.svc.Create(ctx, ports.CreateInput{Items: items})
	if err != nil {
		return nil, err
	}

	r := toResponse(*s)
	return &r, nil
}

func (h *SaleHandler) ListSales(ctx context.Context, req ListSalesRequest) (responses.PageResponse[SaleResponse], error) {
	page, err := h.svc.List(ctx, domain.FindAllQuery{Params: req.Params})
	if err != nil {
		return responses.PageResponse[SaleResponse]{}, err
	}

	return responses.MapPage(page, toResponse), nil
}

func (h *SaleHandler) GetSale(ctx context.Context, id int64) (*SaleResponse, error) {
	s, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	r := toResponse(*s)
	return &r, nil
}
