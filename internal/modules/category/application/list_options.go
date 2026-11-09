package application

import (
	"fmt"
	domain "inventario/internal/modules/category/domain"
)

type ListOptionsUseCase struct {
	repo domain.Repository
}

func NewListOptionsUseCase(repo domain.Repository) *ListOptionsUseCase {
	return &ListOptionsUseCase{repo: repo}
}

func (uc *ListOptionsUseCase) Execute() ([]domain.Option, error) {
	options, err := uc.repo.FindOptions()
	if err != nil {
		return nil, fmt.Errorf("listar opciones de categorías: %w", err)
	}
	return options, nil
}
