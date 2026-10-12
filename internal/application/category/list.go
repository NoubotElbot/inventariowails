package category

import (
	"fmt"
	domain "inventario/internal/domain/category"
)

type ListUseCase struct {
	repo domain.Repository
}

func NewListUseCase(repo domain.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute() ([]domain.Category, error) {
	categories, err := uc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listar categorías: %w", err)
	}
	return categories, nil
}
