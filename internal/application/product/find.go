package product

import (
	"fmt"
	domain "inventario/internal/domain/product"
)

type FindUseCase struct {
	repo domain.Repository
}

func NewFindUseCase(repo domain.Repository) *FindUseCase {
	return &FindUseCase{repo: repo}
}

func (uc *FindUseCase) Execute(id int64) (*domain.Product, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("buscar producto: %w", err)
	}
	return p, nil
}
