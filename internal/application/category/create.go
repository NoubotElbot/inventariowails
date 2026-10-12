package category

import (
	"fmt"
	domain "inventario/internal/domain/category"
	"time"
)

type CreateInput struct {
	Name string
}

type CreateUseCase struct {
	repo domain.Repository
}

func NewCreateUseCase(repo domain.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(input CreateInput) (*domain.Category, error) {
	c := &domain.Category{
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(c); err != nil {
		return nil, fmt.Errorf("crear categoría: %w", err)
	}

	return c, nil
}
