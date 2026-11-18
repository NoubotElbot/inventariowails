// Package services implements the category use cases (the inbound port
// ports.CategoryService) on top of the outbound port ports.CategoryRepository.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
	"inventario/internal/platform/pagination"
)

// CategoryService orchestrates category use cases. It depends only on the
// repository port, never on a concrete database.
type CategoryService struct {
	repo ports.CategoryRepository
}

var _ ports.CategoryService = (*CategoryService)(nil)

func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, input ports.CreateInput) (*domain.Category, error) {
	c := &domain.Category{
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("crear categoría: %w", err)
	}

	return c, nil
}

func (s *CategoryService) Get(ctx context.Context, id int64) (*domain.Category, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar categoría: %w", err)
	}
	return c, nil
}

func (s *CategoryService) List(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	result, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return pagination.Result[domain.Category]{}, fmt.Errorf("listar categorías: %w", err)
	}
	return result, nil
}

func (s *CategoryService) ListOptions(ctx context.Context) ([]domain.Option, error) {
	options, err := s.repo.FindOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("listar opciones de categorías: %w", err)
	}
	return options, nil
}

func (s *CategoryService) Update(ctx context.Context, input ports.UpdateInput) (*domain.Category, error) {
	c, err := s.repo.FindByID(ctx, input.ID)
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

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("actualizar categoría: %w", err)
	}

	return c, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("buscar categoría: %w", err)
	}

	if !c.IsActive() {
		return domain.ErrAlreadyDeleted
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("eliminar categoría: %w", err)
	}

	return nil
}

func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("buscar categoría: %w", err)
	}

	if c.IsActive() {
		return domain.ErrNotDeleted
	}

	if err := s.repo.Restore(ctx, id); err != nil {
		return fmt.Errorf("restaurar categoría: %w", err)
	}

	return nil
}
