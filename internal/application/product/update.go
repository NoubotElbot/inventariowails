package product

import (
	"errors"
	"fmt"
	domain "inventario/internal/domain/product"
	"time"
)

type UpdateInput struct {
	ID         int64
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      float64
	MinStock   int
}

type UpdateUseCase struct {
	repo domain.Repository
}

func NewUpdateUseCase(repo domain.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

func (uc *UpdateUseCase) Execute(input UpdateInput) (*domain.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("buscar producto: %w", err)
	}

	// Verificar duplicado de código solo si cambió
	if p.Code != input.Code {
		_, err := uc.repo.FindByCode(input.Code)
		if err == nil {
			return nil, domain.ErrCodeDuplicate
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("verificar código: %w", err)
		}
	}

	p.CategoryID = input.CategoryID
	p.Name = input.Name
	p.Code = input.Code
	p.Barcode = input.Barcode
	p.Price = input.Price
	p.MinStock = input.MinStock
	p.UpdatedAt = time.Now()

	if err := p.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(p); err != nil {
		return nil, fmt.Errorf("actualizar producto: %w", err)
	}

	return p, nil
}
