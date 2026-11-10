package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "inventario/internal/modules/product/domain"
	"inventario/internal/platform/pagination"
	"inventario/internal/platform/sqlite"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.Repository {
	return &productRepo{db: db}
}

// --- modelo de persistencia ---

type productRow struct {
	ID         int64     `db:"id"`
	CategoryID *int64    `db:"category_id"`
	Name       string    `db:"name"`
	Code       string    `db:"code"`
	Barcode    string    `db:"barcode"`
	Price      int64     `db:"price"`
	Stock      int       `db:"stock"`
	MinStock   int       `db:"min_stock"`
	Active     bool      `db:"active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (r productRow) toDomain() *domain.Product {
	return &domain.Product{
		ID:         r.ID,
		CategoryID: r.CategoryID,
		Name:       r.Name,
		Code:       r.Code,
		Barcode:    r.Barcode,
		Price:      r.Price,
		Stock:      r.Stock,
		MinStock:   r.MinStock,
		Active:     r.Active,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// --- implementación ---

func (r *productRepo) FindAll(q domain.FindAllQuery) (pagination.Result[domain.Product], error) {
	p := q.Params.Normalize()

	const cols = `id, category_id, name, code, barcode, price, stock, min_stock, active, created_at, updated_at`

	var wb sqlite.WhereBuilder
	if q.Search != "" {
		like := "%" + q.Search + "%"
		wb.And(`(name LIKE ? OR code LIKE ? OR barcode LIKE ?)`, like, like, like)
	}
	if len(q.CategoryIDs) > 0 {
		if err := wb.In("category_id", q.CategoryIDs); err != nil {
			return pagination.Result[domain.Product]{}, fmt.Errorf("product.FindAll build in: %w", err)
		}
	}

	var total int64
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM products "+wb.SQL(), wb.Args()...); err != nil {
		return pagination.Result[domain.Product]{}, fmt.Errorf("product.FindAll count: %w", err)
	}

	var rows []productRow
	if err := r.db.Select(&rows,
		"SELECT "+cols+" FROM products "+wb.SQL()+" ORDER BY name ASC LIMIT ? OFFSET ?",
		append(wb.Args(), p.Limit(), p.Offset())...); err != nil {
		return pagination.Result[domain.Product]{}, fmt.Errorf("product.FindAll: %w", err)
	}

	products := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, *row.toDomain())
	}

	return pagination.NewResult(products, total, p), nil
}

func (r *productRepo) FindByID(id int64) (*domain.Product, error) {
	var row productRow

	query := `SELECT id, category_id, name, code, barcode, price,
	                 stock, min_stock, active, created_at, updated_at
	          FROM products
	          WHERE id = ?`

	if err := r.db.Get(&row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("product.FindByID: %w", err)
	}

	return row.toDomain(), nil
}

func (r *productRepo) FindByCode(code string) (*domain.Product, error) {
	var row productRow

	query := `SELECT id, category_id, name, code, barcode, price,
	                 stock, min_stock, active, created_at, updated_at
	          FROM products
	          WHERE code = ?`

	if err := r.db.Get(&row, query, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("product.FindByCode: %w", err)
	}

	return row.toDomain(), nil
}

func (r *productRepo) Save(p *domain.Product) error {
	query := `INSERT INTO products
	            (category_id, name, code, barcode, price, stock, min_stock, active, created_at, updated_at)
	          VALUES
	            (:category_id, :name, :code, :barcode, :price, :stock, :min_stock, :active, :created_at, :updated_at)`

	row := productRow{
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

	result, err := r.db.NamedExec(query, row)
	if err != nil {
		if sqlite.IsUniqueConstraint(err) {
			return domain.ErrCodeDuplicate
		}
		return fmt.Errorf("product.Save: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("product.Save last id: %w", err)
	}

	p.ID = id
	return nil
}

func (r *productRepo) Update(p *domain.Product) error {
	query := `UPDATE products
	          SET category_id = :category_id,
	              name        = :name,
	              code        = :code,
	              barcode     = :barcode,
	              price       = :price,
	              min_stock   = :min_stock,
	              active      = :active,
	              updated_at  = :updated_at
	          WHERE id = :id`

	row := productRow{
		ID:         p.ID,
		CategoryID: p.CategoryID,
		Name:       p.Name,
		Code:       p.Code,
		Barcode:    p.Barcode,
		Price:      p.Price,
		MinStock:   p.MinStock,
		Active:     p.Active,
		UpdatedAt:  p.UpdatedAt,
	}

	result, err := r.db.NamedExec(query, row)
	if err != nil {
		if sqlite.IsUniqueConstraint(err) {
			return domain.ErrCodeDuplicate
		}
		return fmt.Errorf("product.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("product.Update rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *productRepo) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("product.Delete: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("product.Delete rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *productRepo) UpdateStock(id int64, stock int) error {
	query := `UPDATE products
	          SET stock = ?, updated_at = ?
	          WHERE id = ?`

	result, err := r.db.Exec(query, stock, time.Now(), id)
	if err != nil {
		return fmt.Errorf("product.UpdateStock: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("product.UpdateStock rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
