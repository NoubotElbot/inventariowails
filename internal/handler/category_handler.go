package handler

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"inventario/internal/domain/pagination"
)

// --- DTOs ---

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategoryOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

// --- Handler ---

type CategoryHandler struct {
	list        *appCategory.ListUseCase
	listOptions *appCategory.ListOptionsUseCase
	find        *appCategory.FindUseCase
	create      *appCategory.CreateUseCase
	update      *appCategory.UpdateUseCase
	delete      *appCategory.DeleteUseCase
}

func NewCategoryHandler(
	list *appCategory.ListUseCase,
	listOptions *appCategory.ListOptionsUseCase,
	find *appCategory.FindUseCase,
	create *appCategory.CreateUseCase,
	update *appCategory.UpdateUseCase,
	delete *appCategory.DeleteUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		list:        list,
		listOptions: listOptions,
		find:        find,
		create:      create,
		update:      update,
		delete:      delete,
	}
}

// --- métodos expuestos a Wails ---

func (h *CategoryHandler) ListCategories(req pagination.Params) (PageResponse[CategoryResponse], error) {
	page, err := h.list.Execute(req)
	if err != nil {
		return PageResponse[CategoryResponse]{}, err
	}

	return mapPage(page, toResponse), nil
}

func (h *CategoryHandler) ListCategoryOptions() ([]CategoryOptionResponse, error) {
	options, err := h.listOptions.Execute()
	if err != nil {
		return nil, err
	}

	out := make([]CategoryOptionResponse, 0, len(options))
	for _, o := range options {
		out = append(out, CategoryOptionResponse{ID: o.ID, Name: o.Name})
	}
	return out, nil
}

func (h *CategoryHandler) GetCategory(id int64) (*CategoryResponse, error) {
	c, err := h.find.Execute(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) CreateCategory(req CreateCategoryRequest) (*CategoryResponse, error) {
	c, err := h.create.Execute(appCategory.CreateInput{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) UpdateCategory(id int64, req UpdateCategoryRequest) (*CategoryResponse, error) {
	c, err := h.update.Execute(appCategory.UpdateInput{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) DeleteCategory(id int64) error {
	return h.delete.Execute(id)
}

// --- helpers privados ---

func toResponse(c domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
