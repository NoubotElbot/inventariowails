package application

import (
	"fmt"
	domain "inventario/internal/modules/product/domain"
	"inventario/internal/platform/pagination"
)

type ListUseCase struct {
	repo domain.Repository
}

func NewListUseCase(repo domain.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(q domain.FindAllQuery) (pagination.Result[domain.Product], error) {
	page, err := uc.repo.FindAll(q)
	if err != nil {
		return pagination.Result[domain.Product]{}, fmt.Errorf("listar productos: %w", err)
	}
	return page, nil
}
