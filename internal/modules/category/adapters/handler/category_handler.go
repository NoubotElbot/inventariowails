// Package handler is the Wails driving adapter for the category module.
// It translates frontend requests into calls on the ports.CategoryService
// inbound port and maps domain results into response DTOs.
package handler

import (
	"context"
	"errors"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
	"inventario/internal/platform/responses"
)

type CategoryHandler struct {
	svc ports.CategoryService
}

func NewCategoryHandler(svc ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) ListCategories(ctx context.Context, req ListCategoriesRequest) (responses.PageResponse[CategoryResponse], error) {
	page, err := h.svc.List(ctx, domain.FindAllQuery{
		Search:       req.Search,
		OnlyInactive: req.OnlyInactive,
		Params:       req.Params,
	})
	if err != nil {
		return responses.PageResponse[CategoryResponse]{}, err
	}

	return responses.MapPage(page, toResponse), nil
}

func (h *CategoryHandler) ListCategoryOptions(ctx context.Context) ([]CategoryOptionResponse, error) {
	options, err := h.svc.ListOptions(ctx)
	if err != nil {
		return nil, err
	}

	return responses.MapSlice(options, func(o domain.Option) CategoryOptionResponse {
		return CategoryOptionResponse{ID: o.ID, Name: o.Name}
	}), nil
}

func (h *CategoryHandler) GetCategory(ctx context.Context, id int64) (*CategoryResponse, error) {
	c, err := h.svc.Get(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error) {
	c, err := h.svc.Create(ctx, ports.CreateInput{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) UpdateCategory(ctx context.Context, id int64, req UpdateCategoryRequest) (*CategoryResponse, error) {
	c, err := h.svc.Update(ctx, ports.UpdateInput{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, id int64) error {
	return h.svc.Delete(ctx, id)
}

func (h *CategoryHandler) RestoreCategory(ctx context.Context, id int64) error {
	return h.svc.Restore(ctx, id)
}
