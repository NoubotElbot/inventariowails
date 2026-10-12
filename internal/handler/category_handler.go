package handler

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
)

// --- DTOs ---

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

// --- Handler ---

type CategoryHandler struct {
	list   *appCategory.ListUseCase
	find   *appCategory.FindUseCase
	create *appCategory.CreateUseCase
	update *appCategory.UpdateUseCase
	delete *appCategory.DeleteUseCase
}

func NewCategoryHandler(
	list *appCategory.ListUseCase,
	find *appCategory.FindUseCase,
	create *appCategory.CreateUseCase,
	update *appCategory.UpdateUseCase,
	delete *appCategory.DeleteUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		list:   list,
		find:   find,
		create: create,
		update: update,
		delete: delete,
	}
}

// --- métodos expuestos a Wails ---

func (h *CategoryHandler) ListCategories() ([]CategoryResponse, error) {
	categories, err := h.list.Execute()
	if err != nil {
		return nil, err
	}

	response := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		response = append(response, toResponse(c))
	}

	return response, nil
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
		return nil, mapError(err)
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
		return nil, mapError(err)
	}

	r := toResponse(*c)
	return &r, nil
}

func (h *CategoryHandler) DeleteCategory(id int64) error {
	return mapError(h.delete.Execute(id))
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

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return errors.New("categoría no encontrada")
	case errors.Is(err, domain.ErrNameRequired):
		return errors.New("el nombre es requerido")
	case errors.Is(err, domain.ErrNameDuplicate):
		return errors.New("ya existe una categoría con ese nombre")
	default:
		return errors.New("error interno del servidor")
	}
}
