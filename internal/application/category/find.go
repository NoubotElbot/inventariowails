package category

import (
	"fmt"
	domain "inventario/internal/domain/category"
)

type FindUseCase struct {
	repo domain.Repository
}

func NewFindUseCase(repo domain.Repository) *FindUseCase {
	return &FindUseCase{repo: repo}
}

func (uc *FindUseCase) Execute(id int64) (*domain.Category, error) {
	c, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("buscar categoría: %w", err)
	}
	return c, nil
}
