package product

import (
	"errors"
	"fmt"
	domain "inventario/internal/domain/product"
	"time"
)

type CreateInput struct {
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      float64
	Stock      int
	MinStock   int
}

type CreateUseCase struct {
	repo domain.Repository
}

func NewCreateUseCase(repo domain.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(input CreateInput) (*domain.Product, error) {
	_, err := uc.repo.FindByCode(input.Code)
	if err == nil {
		return nil, domain.ErrCodeDuplicate
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("verificar código: %w", err)
	}

	p := &domain.Product{
		CategoryID: input.CategoryID,
		Name:       input.Name,
		Code:       input.Code,
		Barcode:    input.Barcode,
		Price:      input.Price,
		Stock:      input.Stock,
		MinStock:   input.MinStock,
		Active:     true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(p); err != nil {
		return nil, fmt.Errorf("crear producto: %w", err)
	}

	return p, nil
}
