package product

import (
	"fmt"
	domain "inventario/internal/domain/product"
)

type ListUseCase struct {
	repo domain.Repository
}

func NewListUseCase(repo domain.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute() ([]domain.Product, error) {
	products, err := uc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listar productos: %w", err)
	}
	return products, nil
}
