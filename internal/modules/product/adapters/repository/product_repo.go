// Package repository is the GORM driven adapter implementing
// ports.ProductRepository for the product module.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"

	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

var _ ports.ProductRepository = (*productRepo)(nil)

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepo{db: db}
}

// --- modelo de persistencia ---

type productModel struct {
	ID         int64 `gorm:"primaryKey"`
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      int64
	Stock      int
	MinStock   int
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (productModel) TableName() string { return "products" }

func (m productModel) toDomain() *domain.Product {
	return &domain.Product{
		ID:         m.ID,
		CategoryID: m.CategoryID,
		Name:       m.Name,
		Code:       m.Code,
		Barcode:    m.Barcode,
		Price:      m.Price,
		Stock:      m.Stock,
		MinStock:   m.MinStock,
		Active:     m.Active,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func fromDomain(p *domain.Product) productModel {
	return productModel{
		ID:         p.ID,
		CategoryID: p.CategoryID,
		Name:       p.Name,
		Code:       p.Code,
		Barcode:    p.Barcode,
		Price:      p.Price,
		Stock:      p.Stock,
		MinStock:   p.MinStock,
		Active:     p.Active,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// updatableColumns are the columns Update may change. Stock is excluded (it
// only changes through UpdateStock) and so is created_at.
var updatableColumns = []string{
	"category_id", "name", "code", "barcode", "price", "min_stock", "active", "updated_at",
}

// --- implementación ---

func (r *productRepo) FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Product], error) {
	query := gormdb.Conn(ctx, r.db).
		Model(&productModel{}).
		Order("name ASC").
		Scopes(findAllFilters(q))

	res, err := pagination.Paginate[domain.Product](query, q.Params)
	if err != nil {
		return res, fmt.Errorf("product.FindAll: %w", err)
	}
	return res, nil
}

func (r *productRepo) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	var row productModel
	if err := gormdb.Conn(ctx, r.db).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("product.FindByID: %w", err)
	}
	return row.toDomain(), nil
}

func (r *productRepo) FindByCode(ctx context.Context, code string) (*domain.Product, error) {
	var row productModel
	if err := gormdb.Conn(ctx, r.db).Where("code = ?", code).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("product.FindByCode: %w", err)
	}
	return row.toDomain(), nil
}

func (r *productRepo) Save(ctx context.Context, p *domain.Product) error {
	row := fromDomain(p)
	row.ID = 0

	if err := gormdb.Conn(ctx, r.db).Create(&row).Error; err != nil {
		if gormdb.IsUniqueConstraint(err) {
			return domain.ErrCodeDuplicate
		}
		return fmt.Errorf("product.Save: %w", err)
	}

	p.ID = row.ID
	return nil
}

func (r *productRepo) Update(ctx context.Context, p *domain.Product) error {
	row := fromDomain(p)

	result := gormdb.Conn(ctx, r.db).
		Model(&productModel{ID: p.ID}).
		Select(updatableColumns).
		Updates(&row)
	if result.Error != nil {
		if gormdb.IsUniqueConstraint(result.Error) {
			return domain.ErrCodeDuplicate
		}
		return fmt.Errorf("product.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *productRepo) Delete(ctx context.Context, id int64) error {
	result := gormdb.Conn(ctx, r.db).Delete(&productModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("product.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepo) UpdateStock(ctx context.Context, id int64, stock int) error {
	result := gormdb.Conn(ctx, r.db).
		Model(&productModel{}).
		Where("id = ?", id).
		Updates(map[string]any{"stock": stock, "updated_at": time.Now()})
	if result.Error != nil {
		return fmt.Errorf("product.UpdateStock: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
