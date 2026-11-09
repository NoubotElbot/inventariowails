package application

import (
	"errors"
	"fmt"
	domain "inventario/internal/modules/category/domain"
)

type DeleteUseCase struct {
	repo domain.Repository
}

func NewDeleteUseCase(repo domain.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(id int64) error {
	c, err := uc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("buscar categoría: %w", err)
	}

	if !c.IsActive() {
		return domain.ErrAlreadyDeleted
	}

	if err := uc.repo.SoftDelete(id); err != nil {
		return fmt.Errorf("eliminar categoría: %w", err)
	}

	return nil
}
