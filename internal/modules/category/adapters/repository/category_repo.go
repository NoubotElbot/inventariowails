// Package repository is the GORM driven adapter implementing
// ports.CategoryRepository for the category module.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"

	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

var _ ports.CategoryRepository = (*categoryRepo)(nil)

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &categoryRepo{db: db}
}

// --- modelo de persistencia ---

// categoryModel maps the categories table. DeletedAt is a plain *time.Time
// (not gorm.DeletedAt) on purpose: soft-delete filtering is explicit in each
// query, so FindByID can still return deleted rows and FindAll can list only
// inactive ones.
type categoryModel struct {
	ID        int64 `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (categoryModel) TableName() string { return "categories" }

func (m categoryModel) toDomain() *domain.Category {
	return &domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func fromDomain(c *domain.Category) categoryModel {
	return categoryModel{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

// --- implementación ---

func (r *categoryRepo) FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	query := r.db.Model(&categoryModel{}).
		Order("id DESC").
		Scopes(findAllFilters(q))

	return pagination.Paginate[domain.Category](query, q.Params)
}

func (r *categoryRepo) FindOptions(ctx context.Context) ([]domain.Option, error) {
	var rows []categoryModel
	if err := gormdb.Conn(ctx, r.db).
		Select("id", "name").
		Order("name ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("category.FindOptions: %w", err)
	}

	options := make([]domain.Option, 0, len(rows))
	for _, row := range rows {
		options = append(options, domain.Option{ID: row.ID, Name: row.Name})
	}
	return options, nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	var row categoryModel
	if err := gormdb.Conn(ctx, r.db).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("category.FindByID: %w", err)
	}
	return row.toDomain(), nil
}

func (r *categoryRepo) Save(ctx context.Context, c *domain.Category) error {
	row := fromDomain(c)

	if err := gormdb.Conn(ctx, r.db).Create(&row).Error; err != nil {
		if gormdb.IsUniqueConstraint(err) {
			return domain.ErrNameDuplicate
		}
		return fmt.Errorf("category.Save: %w", err)
	}

	c.ID = row.ID
	return nil
}

func (r *categoryRepo) Update(ctx context.Context, c *domain.Category) error {
	row := fromDomain(c)
	result := gormdb.Conn(ctx, r.db).
		Model(&row).
		Select("name", "updated_at").
		Updates(&row)
	if result.Error != nil {
		if gormdb.IsUniqueConstraint(result.Error) {
			return domain.ErrNameDuplicate
		}
		return fmt.Errorf("category.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *categoryRepo) SoftDelete(ctx context.Context, id int64) error {
	result := gormdb.Conn(ctx, r.db).
		Model(&categoryModel{}).
		Scopes(active).
		Where("id = ?", id).
		UpdateColumn("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("category.SoftDelete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *categoryRepo) Restore(ctx context.Context, id int64) error {
	result := gormdb.Conn(ctx, r.db).
		Model(&categoryModel{}).
		Scopes(deleted).
		Where("id = ?", id).
		UpdateColumn("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("category.Restore: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
