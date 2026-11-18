// Package services implements the product use cases (the inbound port
// ports.ProductService) on top of the outbound port ports.ProductRepository.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/platform/pagination"
)

// ProductService orchestrates product use cases. It depends only on the
// repository port, never on a concrete database.
type ProductService struct {
	repo ports.ProductRepository
}

var _ ports.ProductService = (*ProductService)(nil)

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, input ports.CreateInput) (*domain.Product, error) {
	if err := s.ensureCodeAvailable(ctx, input.Code); err != nil {
		return nil, err
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

	if err := s.repo.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("crear producto: %w", err)
	}

	return p, nil
}

func (s *ProductService) Get(ctx context.Context, id int64) (*domain.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar producto: %w", err)
	}
	return p, nil
}

func (s *ProductService) List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Product], error) {
	result, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return pagination.Result[domain.Product]{}, fmt.Errorf("listar productos: %w", err)
	}
	return result, nil
}

func (s *ProductService) Update(ctx context.Context, input ports.UpdateInput) (*domain.Product, error) {
	p, err := s.findExisting(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Verificar duplicado de código solo si cambió
	if p.Code != input.Code {
		if err := s.ensureCodeAvailable(ctx, input.Code); err != nil {
			return nil, err
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

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("actualizar producto: %w", err)
	}

	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if _, err := s.findExisting(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("eliminar producto: %w", err)
	}

	return nil
}

func (s *ProductService) ToggleActive(ctx context.Context, id int64) (*domain.Product, error) {
	p, err := s.findExisting(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Active = !p.Active
	p.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("actualizar producto: %w", err)
	}

	return p, nil
}

// findExisting loads a product, returning domain.ErrNotFound unwrapped when
// it does not exist.
func (s *ProductService) findExisting(ctx context.Context, id int64) (*domain.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("buscar producto: %w", err)
	}
	return p, nil
}

// ensureCodeAvailable returns domain.ErrCodeDuplicate if another product
// already uses code.
func (s *ProductService) ensureCodeAvailable(ctx context.Context, code string) error {
	_, err := s.repo.FindByCode(ctx, code)
	if err == nil {
		return domain.ErrCodeDuplicate
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("verificar código: %w", err)
	}
	return nil
}
