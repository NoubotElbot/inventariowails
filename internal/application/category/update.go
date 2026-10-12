package category

import (
	"errors"
	"fmt"
	domain "inventario/internal/domain/category"
	"time"
)

type UpdateInput struct {
	ID   int64
	Name string
}

type UpdateUseCase struct {
	repo domain.Repository
}

func NewUpdateUseCase(repo domain.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

func (uc *UpdateUseCase) Execute(input UpdateInput) (*domain.Category, error) {
	c, err := uc.repo.FindByID(input.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("buscar categoría: %w", err)
	}

	c.Name = input.Name
	c.UpdatedAt = time.Now()

	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(c); err != nil {
		return nil, fmt.Errorf("actualizar categoría: %w", err)
	}

	return c, nil
}
