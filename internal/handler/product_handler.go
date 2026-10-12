package handler

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
)

// --- DTOs ---

type ProductResponse struct {
	ID         int64   `json:"id"`
	CategoryID *int64  `json:"category_id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	MinStock   int     `json:"min_stock"`
	Active     bool    `json:"active"`
	LowStock   bool    `json:"low_stock"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type CreateProductRequest struct {
	CategoryID *int64  `json:"category_id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	MinStock   int     `json:"min_stock"`
}

type UpdateProductRequest struct {
	CategoryID *int64  `json:"category_id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
	MinStock   int     `json:"min_stock"`
}

// --- Handler ---

type ProductHandler struct {
	list         *appProduct.ListUseCase
	find         *appProduct.FindUseCase
	create       *appProduct.CreateUseCase
	update       *appProduct.UpdateUseCase
	delete       *appProduct.DeleteUseCase
	toggleActive *appProduct.ToggleActiveUseCase
}

func NewProductHandler(
	list *appProduct.ListUseCase,
	find *appProduct.FindUseCase,
	create *appProduct.CreateUseCase,
	update *appProduct.UpdateUseCase,
	delete *appProduct.DeleteUseCase,
	toggleActive *appProduct.ToggleActiveUseCase,
) *ProductHandler {
	return &ProductHandler{
		list:         list,
		find:         find,
		create:       create,
		update:       update,
		delete:       delete,
		toggleActive: toggleActive,
	}
}

// --- métodos expuestos a Wails ---

func (h *ProductHandler) ListProducts() ([]ProductResponse, error) {
	products, err := h.list.Execute()
	if err != nil {
		return nil, err
	}

	response := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		response = append(response, toProductResponse(p))
	}

	return response, nil
}

func (h *ProductHandler) GetProduct(id int64) (*ProductResponse, error) {
	p, err := h.find.Execute(id)
	if err != nil {
		return nil, mapProductError(err)
	}

	r := toProductResponse(*p)
	return &r, nil
}

func (h *ProductHandler) CreateProduct(req CreateProductRequest) (*ProductResponse, error) {
	p, err := h.create.Execute(appProduct.CreateInput{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Code:       req.Code,
		Barcode:    req.Barcode,
		Price:      req.Price,
		Stock:      req.Stock,
		MinStock:   req.MinStock,
	})
	if err != nil {
		return nil, mapProductError(err)
	}

	r := toProductResponse(*p)
	return &r, nil
}

func (h *ProductHandler) UpdateProduct(id int64, req UpdateProductRequest) (*ProductResponse, error) {
	p, err := h.update.Execute(appProduct.UpdateInput{
		ID:         id,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Code:       req.Code,
		Barcode:    req.Barcode,
		Price:      req.Price,
		MinStock:   req.MinStock,
	})
	if err != nil {
		return nil, mapProductError(err)
	}

	r := toProductResponse(*p)
	return &r, nil
}

func (h *ProductHandler) DeleteProduct(id int64) error {
	return mapProductError(h.delete.Execute(id))
}

func (h *ProductHandler) ToggleActiveProduct(id int64) (*ProductResponse, error) {
	p, err := h.toggleActive.Execute(id)
	if err != nil {
		return nil, mapProductError(err)
	}

	r := toProductResponse(*p)
	return &r, nil
}

// --- helpers privados ---

func toProductResponse(p domain.Product) ProductResponse {
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
		CreatedAt:  p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func mapProductError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return errors.New("producto no encontrado")
	case errors.Is(err, domain.ErrNameRequired):
		return errors.New("el nombre es requerido")
	case errors.Is(err, domain.ErrCodeRequired):
		return errors.New("el código es requerido")
	case errors.Is(err, domain.ErrCodeDuplicate):
		return errors.New("ya existe un producto con ese código")
	case errors.Is(err, domain.ErrPriceInvalid):
		return errors.New("el precio debe ser mayor a cero")
	case errors.Is(err, domain.ErrStockNegative):
		return errors.New("el stock no puede ser negativo")
	case errors.Is(err, domain.ErrInsufficientStock):
		return errors.New("stock insuficiente")
	default:
		return errors.New("error interno del servidor")
	}
}
