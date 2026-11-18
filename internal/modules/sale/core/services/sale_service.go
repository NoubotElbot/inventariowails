// Package services implements the sale use cases (the inbound port
// ports.SaleService) on top of the outbound ports for sales, products and
// stock movements.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	stockdomain "inventario/internal/modules/stock/core/domain"
	stockports "inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/pagination"
	"inventario/internal/platform/uow"
)

// SaleService orchestrates sale use cases. It depends only on ports, never
// on a concrete database or on other modules' implementations.
type SaleService struct {
	uow      uow.UnitOfWork
	repo     ports.SaleRepository
	products ports.ProductRepository
	stock    ports.StockMovementRegistrar
}

var _ ports.SaleService = (*SaleService)(nil)

func NewSaleService(
	u uow.UnitOfWork,
	repo ports.SaleRepository,
	products ports.ProductRepository,
	stock ports.StockMovementRegistrar,
) *SaleService {
	return &SaleService{uow: u, repo: repo, products: products, stock: stock}
}

// Create prices each item from the current product data, then saves the sale
// and takes the stock out of every product in a single transaction.
func (s *SaleService) Create(ctx context.Context, input ports.CreateInput) (*domain.Sale, error) {
	if len(input.Items) == 0 {
		return nil, domain.ErrNoItems
	}

	sale := &domain.Sale{CreatedAt: time.Now()}

	for _, in := range input.Items {
		product, err := s.findSellableProduct(ctx, in)
		if err != nil {
			return nil, err
		}

		subtotal := product.Price * int64(in.Quantity)
		sale.Items = append(sale.Items, domain.SaleItem{
			ProductID: product.ID,
			Quantity:  in.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
			CreatedAt: time.Now(),
		})
		sale.Total += subtotal
	}

	if err := sale.Validate(); err != nil {
		return nil, err
	}

	// Guardar la venta y descontar el stock de cada producto en una sola
	// transacción. El descuento lo hace el módulo de stock, que actualiza el
	// stock del producto y registra el movimiento dentro de esta misma
	// transacción.
	err := s.uow.Do(ctx, func(ctx context.Context) error {
		if err := s.repo.Save(ctx, sale); err != nil {
			return fmt.Errorf("guardar venta: %w", err)
		}

		for _, item := range sale.Items {
			if _, err := s.stock.RegisterMovement(ctx, stockports.RegisterMovementInput{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Type:      stockdomain.MovementSale,
				Notes:     fmt.Sprintf("venta #%d", sale.ID),
			}); err != nil {
				return fmt.Errorf("descontar stock producto %d: %w", item.ProductID, err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *SaleService) Get(ctx context.Context, id int64) (*domain.Sale, error) {
	sale, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar venta: %w", err)
	}
	return sale, nil
}

func (s *SaleService) List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Sale], error) {
	page, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return pagination.Result[domain.Sale]{}, fmt.Errorf("listar ventas: %w", err)
	}
	return page, nil
}

// findSellableProduct loads the item's product and checks it can be sold:
// it exists, is active and has enough stock. The stock module re-checks the
// stock inside the transaction.
func (s *SaleService) findSellableProduct(ctx context.Context, in ports.SaleItemInput) (*productdomain.Product, error) {
	product, err := s.products.FindByID(ctx, in.ProductID)
	if err != nil {
		if errors.Is(err, productdomain.ErrNotFound) {
			return nil, productdomain.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto %d: %w", in.ProductID, err)
	}

	if !product.Active {
		return nil, fmt.Errorf("producto %d: %w", in.ProductID, domain.ErrProductInactive)
	}

	if product.Stock < in.Quantity {
		return nil, productdomain.ErrInsufficientStock
	}

	return product, nil
}
