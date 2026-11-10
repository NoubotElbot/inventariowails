package application

import (
	"errors"
	"fmt"
	domain "inventario/internal/modules/category/domain"
)

type RestoreUseCase struct {
	repo domain.Repository
}

func NewRestoreUseCase(repo domain.Repository) *RestoreUseCase {
	return &RestoreUseCase{repo: repo}
}

func (uc *RestoreUseCase) Execute(id int64) error {
	c, err := uc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("buscar categoría: %w", err)
	}

	if c.IsActive() {
		return domain.ErrNotDeleted
	}

	if err := uc.repo.Restore(id); err != nil {
		return fmt.Errorf("restaurar categoría: %w", err)
	}

	return nil
}
