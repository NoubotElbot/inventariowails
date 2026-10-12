package product

import (
	"errors"
	"fmt"
	domain "inventario/internal/domain/product"
	"time"
)

type ToggleActiveUseCase struct {
	repo domain.Repository
}

func NewToggleActiveUseCase(repo domain.Repository) *ToggleActiveUseCase {
	return &ToggleActiveUseCase{repo: repo}
}

func (uc *ToggleActiveUseCase) Execute(id int64) (*domain.Product, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto: %w", err)
	}

	p.Active = !p.Active
	p.UpdatedAt = time.Now()

	if err := uc.repo.Update(p); err != nil {
		return nil, fmt.Errorf("actualizar producto: %w", err)
	}

	return p, nil
}
