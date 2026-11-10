package application

import (
	"errors"
	"fmt"
	domainProduct "inventario/internal/modules/product/domain"
	domain "inventario/internal/modules/stock/domain"
)

type ListMovementsUseCase struct {
	stockRepo   domain.Repository
	productRepo domainProduct.Repository
}

func NewListMovementsUseCase(
	stockRepo domain.Repository,
	productRepo domainProduct.Repository,
) *ListMovementsUseCase {
	return &ListMovementsUseCase{
		stockRepo:   stockRepo,
		productRepo: productRepo,
	}
}

func (uc *ListMovementsUseCase) Execute(productID int64) ([]domain.Movement, error) {
	_, err := uc.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, domainProduct.ErrNotFound) {
			return nil, domainProduct.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto: %w", err)
	}

	movements, err := uc.stockRepo.FindByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("listar movimientos: %w", err)
	}

	return movements, nil
}
