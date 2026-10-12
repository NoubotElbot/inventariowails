package stock

import (
	"errors"
	"fmt"
	domainProduct "inventario/internal/domain/product"
	domain "inventario/internal/domain/stock"
	"time"
)

type AdjustInput struct {
	ProductID int64
	Quantity  int
	Notes     string
}

type AdjustUseCase struct {
	stockRepo   domain.Repository
	productRepo domainProduct.Repository
}

func NewAdjustUseCase(
	stockRepo domain.Repository,
	productRepo domainProduct.Repository,
) *AdjustUseCase {
	return &AdjustUseCase{
		stockRepo:   stockRepo,
		productRepo: productRepo,
	}
}

func (uc *AdjustUseCase) Execute(input AdjustInput) (*domain.Movement, error) {
	product, err := uc.productRepo.FindByID(input.ProductID)
	if err != nil {
		if errors.Is(err, domainProduct.ErrNotFound) {
			return nil, domainProduct.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto: %w", err)
	}

	m := &domain.Movement{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		MovementType: domain.MovementAdjustment,
		Notes:        input.Notes,
		CreatedAt:    time.Now(),
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	newStock := product.Stock + input.Quantity
	if newStock < 0 {
		return nil, domainProduct.ErrInsufficientStock
	}

	if err := uc.productRepo.UpdateStock(input.ProductID, newStock); err != nil {
		return nil, fmt.Errorf("actualizar stock: %w", err)
	}

	if err := uc.stockRepo.Save(m); err != nil {
		return nil, fmt.Errorf("registrar movimiento: %w", err)
	}

	return m, nil
}
