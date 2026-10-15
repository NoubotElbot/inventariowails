package sale

import (
	"errors"
	"fmt"
	domainProduct "inventario/internal/domain/product"
	domainSale "inventario/internal/domain/sale"
	domainStock "inventario/internal/domain/stock"
	"time"
)

type SaleItemInput struct {
	ProductID int64
	Quantity  int
}

type CreateInput struct {
	Items []SaleItemInput
}

type CreateUseCase struct {
	saleRepo    domainSale.Repository
	productRepo domainProduct.Repository
	stockRepo   domainStock.Repository
}

func NewCreateUseCase(
	saleRepo domainSale.Repository,
	productRepo domainProduct.Repository,
	stockRepo domainStock.Repository,
) *CreateUseCase {
	return &CreateUseCase{
		saleRepo:    saleRepo,
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

func (uc *CreateUseCase) Execute(input CreateInput) (*domainSale.Sale, error) {
	if len(input.Items) == 0 {
		return nil, domainSale.ErrNoItems
	}

	sale := &domainSale.Sale{
		CreatedAt: time.Now(),
	}

	for _, itemInput := range input.Items {
		product, err := uc.productRepo.FindByID(itemInput.ProductID)
		if err != nil {
			if errors.Is(err, domainProduct.ErrNotFound) {
				return nil, domainProduct.ErrNotFound
			}
			return nil, fmt.Errorf("buscar producto %d: %w", itemInput.ProductID, err)
		}

		if !product.Active {
			return nil, fmt.Errorf("producto %d inactivo", itemInput.ProductID)
		}

		if product.Stock < itemInput.Quantity {
			return nil, domainProduct.ErrInsufficientStock
		}

		subtotal := product.Price * float64(itemInput.Quantity)

		sale.Items = append(sale.Items, domainSale.SaleItem{
			ProductID: product.ID,
			Quantity:  itemInput.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
			CreatedAt: time.Now(),
		})

		sale.Total += subtotal
	}

	if err := sale.Validate(); err != nil {
		return nil, err
	}

	if err := uc.saleRepo.Save(sale); err != nil {
		return nil, fmt.Errorf("guardar venta: %w", err)
	}

	// Descontar stock y registrar movimientos
	for _, item := range sale.Items {
		product, err := uc.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("recargar producto %d: %w", item.ProductID, err)
		}

		newStock := product.Stock - item.Quantity
		if err := uc.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
			return nil, fmt.Errorf("actualizar stock %d: %w", item.ProductID, err)
		}

		movement := &domainStock.Movement{
			ProductID:    item.ProductID,
			Quantity:     -item.Quantity,
			MovementType: domainStock.MovementSale,
			Notes:        fmt.Sprintf("venta #%d", sale.ID),
			CreatedAt:    time.Now(),
		}
		if err := uc.stockRepo.Save(movement); err != nil {
			return nil, fmt.Errorf("registrar movimiento de stock: %w", err)
		}
	}

	return sale, nil
}
