package application

import (
	"fmt"
	domain "inventario/internal/modules/category/domain"
	"inventario/internal/platform/pagination"
)

type ListUseCase struct {
	repo domain.Repository
}

func NewListUseCase(repo domain.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	result, err := uc.repo.FindAll(q)
	if err != nil {
		return pagination.Result[domain.Category]{}, fmt.Errorf("listar categorías: %w", err)
	}
	return result, nil
}
