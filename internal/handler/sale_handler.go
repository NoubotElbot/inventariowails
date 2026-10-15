package handler

import (
	"errors"
	appSale "inventario/internal/application/sale"
	domainProduct "inventario/internal/domain/product"
	domainSale "inventario/internal/domain/sale"
)

// --- DTOs ---

type SaleItemResponse struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
	CreatedAt string  `json:"created_at"`
}

type SaleResponse struct {
	ID        int64              `json:"id"`
	Total     float64            `json:"total"`
	Items     []SaleItemResponse `json:"items"`
	CreatedAt string             `json:"created_at"`
}

type SaleItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CreateSaleRequest struct {
	Items []SaleItemRequest `json:"items"`
}

// --- Handler ---

type SaleHandler struct {
	create *appSale.CreateUseCase
	list   *appSale.ListUseCase
	find   *appSale.FindUseCase
}

func NewSaleHandler(
	create *appSale.CreateUseCase,
	list *appSale.ListUseCase,
	find *appSale.FindUseCase,
) *SaleHandler {
	return &SaleHandler{
		create: create,
		list:   list,
		find:   find,
	}
}

// --- métodos expuestos a Wails ---

func (h *SaleHandler) CreateSale(req CreateSaleRequest) (*SaleResponse, error) {
	items := make([]appSale.SaleItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, appSale.SaleItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	s, err := h.create.Execute(appSale.CreateInput{Items: items})
	if err != nil {
		return nil, mapSaleError(err)
	}

	r := toSaleResponse(*s)
	return &r, nil
}

func (h *SaleHandler) ListSales() ([]SaleResponse, error) {
	sales, err := h.list.Execute()
	if err != nil {
		return nil, err
	}

	response := make([]SaleResponse, 0, len(sales))
	for _, s := range sales {
		response = append(response, toSaleResponse(s))
	}

	return response, nil
}

func (h *SaleHandler) GetSale(id int64) (*SaleResponse, error) {
	s, err := h.find.Execute(id)
	if err != nil {
		return nil, mapSaleError(err)
	}

	r := toSaleResponse(*s)
	return &r, nil
}

// --- helpers privados ---

func toSaleResponse(s domainSale.Sale) SaleResponse {
	items := make([]SaleItemResponse, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, SaleItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return SaleResponse{
		ID:        s.ID,
		Total:     s.Total,
		Items:     items,
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func mapSaleError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domainSale.ErrNotFound):
		return errors.New("venta no encontrada")
	case errors.Is(err, domainSale.ErrNoItems):
		return errors.New("la venta debe tener al menos un producto")
	case errors.Is(err, domainSale.ErrQuantityInvalid):
		return errors.New("la cantidad debe ser mayor a cero")
	case errors.Is(err, domainSale.ErrProductInvalid):
		return errors.New("producto inválido en la venta")
	case errors.Is(err, domainProduct.ErrNotFound):
		return errors.New("producto no encontrado")
	case errors.Is(err, domainProduct.ErrInsufficientStock):
		return errors.New("stock insuficiente")
	default:
		return errors.New("error interno del servidor")
	}
}
