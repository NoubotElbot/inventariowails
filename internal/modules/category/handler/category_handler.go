package handler

import (
	"errors"
	appCategory "inventario/internal/modules/category/application"
	domain "inventario/internal/modules/category/domain"
	"inventario/internal/platform/pagination"
	"inventario/internal/platform/responses"
)

// ListCategoriesRequest carries pagination and optional search filter.
type ListCategoriesRequest struct {
	Search       string `json:"search"`
	OnlyInactive bool   `json:"only_inactive"`
	pagination.Params
}

// --- DTOs ---

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
	restore     *appCategory.RestoreUseCase
}

func NewCategoryHandler(
	list *appCategory.ListUseCase,
	listOptions *appCategory.ListOptionsUseCase,
	find *appCategory.FindUseCase,
	create *appCategory.CreateUseCase,
	update *appCategory.UpdateUseCase,
	delete *appCategory.DeleteUseCase,
	restore *appCategory.RestoreUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		list:        list,
		listOptions: listOptions,
		find:        find,
		create:      create,
		update:      update,
		delete:      delete,
		restore:     restore,
	}
}

// --- métodos expuestos a Wails ---

func (h *CategoryHandler) ListCategories(req ListCategoriesRequest) (responses.PageResponse[CategoryResponse], error) {
	page, err := h.list.Execute(domain.FindAllQuery{
		Search:       req.Search,
		OnlyInactive: req.OnlyInactive,
		Params:       req.Params,
	})
	if err != nil {
		return responses.PageResponse[CategoryResponse]{}, err
	}

	return responses.MapPage(page, toResponse), nil
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

func (h *CategoryHandler) RestoreCategory(id int64) error {
	return h.restore.Execute(id)
}

// --- helpers privados ---

func toResponse(c domain.Category) CategoryResponse {
	r := CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if c.DeletedAt != nil {
		s := c.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
		r.DeletedAt = &s
	}
	return r
}
