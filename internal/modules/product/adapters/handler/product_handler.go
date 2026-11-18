// Package handler is the Wails driving adapter for the product module.
// It translates frontend requests into calls on the ports.ProductService
// inbound port and maps domain results into response DTOs.
package handler

import (
	"context"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/platform/responses"
)

type ProductHandler struct {
	svc ports.ProductService
}

func NewProductHandler(svc ports.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) ListProducts(ctx context.Context, req ListProductsRequest) (responses.PageResponse[ProductResponse], error) {
	page, err := h.svc.List(ctx, domain.FindAllQuery{
		Search:      req.Search,
		CategoryIDs: req.CategoryIDs,
		Params:      req.Params,
	})
	if err != nil {
		return responses.PageResponse[ProductResponse]{}, err
	}

	return responses.MapPage(page, toResponse), nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, id int64) (*ProductResponse, error) {
	p, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	r := toResponse(*p)
	return &r, nil
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	p, err := h.svc.Create(ctx, ports.CreateInput{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Code:       req.Code,
		Barcode:    req.Barcode,
		Price:      req.Price,
		Stock:      req.Stock,
		MinStock:   req.MinStock,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*p)
	return &r, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (*ProductResponse, error) {
	p, err := h.svc.Update(ctx, ports.UpdateInput{
		ID:         id,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Code:       req.Code,
		Barcode:    req.Barcode,
		Price:      req.Price,
		MinStock:   req.MinStock,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*p)
	return &r, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, id int64) error {
	return h.svc.Delete(ctx, id)
}

func (h *ProductHandler) ToggleActiveProduct(ctx context.Context, id int64) (*ProductResponse, error) {
	p, err := h.svc.ToggleActive(ctx, id)
	if err != nil {
		return nil, err
	}

	r := toResponse(*p)
	return &r, nil
}
