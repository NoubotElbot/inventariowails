package category

import (
	"fmt"
	domain "inventario/internal/domain/category"
	"inventario/internal/domain/pagination"
)

type ListUseCase struct {
	repo domain.Repository
}

func NewListUseCase(repo domain.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(p pagination.Params) (pagination.Result[domain.Category], error) {
	result, err := uc.repo.FindAll(p.Normalize())
	if err != nil {
		return pagination.Result[domain.Category]{}, fmt.Errorf("listar categorías: %w", err)
	}
	return result, nil
}
