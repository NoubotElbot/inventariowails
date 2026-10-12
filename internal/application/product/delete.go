package product

import (
	"errors"
	"fmt"
	domain "inventario/internal/domain/product"
)

type DeleteUseCase struct {
	repo domain.Repository
}

func NewDeleteUseCase(repo domain.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(id int64) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("buscar producto: %w", err)
	}

	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("eliminar producto: %w", err)
	}

	return nil
}
