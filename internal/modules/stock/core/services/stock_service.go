// Package services implements the stock use cases (the inbound port
// ports.StockService) on top of the outbound ports ports.StockRepository and
// ports.ProductRepository.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/uow"
)

// StockService is the single place where product stock changes: it applies
// movements to the product's stock and records them.
type StockService struct {
	uow      uow.UnitOfWork
	repo     ports.StockRepository
	products ports.ProductRepository
}

var _ ports.StockService = (*StockService)(nil)

func NewStockService(u uow.UnitOfWork, repo ports.StockRepository, products ports.ProductRepository) *StockService {
	return &StockService{uow: u, repo: repo, products: products}
}

// AdjustStock registers a manual adjustment. Only adjustment_in and
// adjustment_out are accepted: sales and returns are registered by their own
// modules through RegisterMovement.
func (s *StockService) AdjustStock(ctx context.Context, input ports.AdjustInput) (*domain.Movement, error) {
	movement := domain.Movement{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		MovementType: input.Type,
	}
	if err := movement.Validate(); err != nil {
		return nil, err
	}
	if !input.Type.IsManualAdjustment() {
		return nil, domain.ErrInvalidMovement
	}

	return s.RegisterMovement(ctx, ports.RegisterMovementInput(input))
}

// RegisterMovement reads the product, updates its stock and records the
// movement in a single transaction, so stock and movement history never
// diverge. Called inside another unit of work (e.g. creating a sale) it joins
// that transaction instead of opening its own.
func (s *StockService) RegisterMovement(ctx context.Context, input ports.RegisterMovementInput) (*domain.Movement, error) {
	movement := &domain.Movement{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		MovementType: input.Type,
		Notes:        input.Notes,
		CreatedAt:    time.Now(),
	}

	if err := movement.Validate(); err != nil {
		return nil, err
	}

	err := s.uow.Do(ctx, func(ctx context.Context) error {
		product, err := s.findProduct(ctx, input.ProductID)
		if err != nil {
			return err
		}

		newStock := product.Stock + movement.Delta()
		if newStock < 0 {
			return productdomain.ErrInsufficientStock
		}

		if err := s.products.UpdateStock(ctx, input.ProductID, newStock); err != nil {
			return fmt.Errorf("actualizar stock: %w", err)
		}

		if err := s.repo.Save(ctx, movement); err != nil {
			return fmt.Errorf("registrar movimiento: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return movement, nil
}

func (s *StockService) ListMovements(ctx context.Context, productID int64) ([]domain.Movement, error) {
	if productID == 0 {
		return nil, domain.ErrProductRequired
	}

	if _, err := s.findProduct(ctx, productID); err != nil {
		return nil, err
	}

	movements, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("listar movimientos: %w", err)
	}

	return movements, nil
}

// findProduct loads a product, returning productdomain.ErrNotFound unwrapped
// when it does not exist.
func (s *StockService) findProduct(ctx context.Context, id int64) (*productdomain.Product, error) {
	p, err := s.products.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, productdomain.ErrNotFound) {
			return nil, productdomain.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto: %w", err)
	}
	return p, nil
}
