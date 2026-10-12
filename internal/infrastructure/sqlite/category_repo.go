package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "inventario/internal/domain/category"

	"github.com/jmoiron/sqlx"
)

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) domain.Repository {
	return &categoryRepo{db: db}
}

// --- modelo de persistencia ---

type categoryRow struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r categoryRow) toDomain() *domain.Category {
	return &domain.Category{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// --- implementación ---

func (r *categoryRepo) FindAll() ([]domain.Category, error) {
	var rows []categoryRow

	query := `SELECT id, name, created_at, updated_at
	          FROM categories
	          ORDER BY name ASC`

	if err := r.db.Select(&rows, query); err != nil {
		return nil, fmt.Errorf("category.FindAll: %w", err)
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, *row.toDomain())
	}

	return categories, nil
}

func (r *categoryRepo) FindByID(id int64) (*domain.Category, error) {
	var row categoryRow

	query := `SELECT id, name, created_at, updated_at
	          FROM categories
	          WHERE id = ?`

	if err := r.db.Get(&row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("category.FindByID: %w", err)
	}

	return row.toDomain(), nil
}

func (r *categoryRepo) Save(c *domain.Category) error {
	query := `INSERT INTO categories (name, created_at, updated_at)
	          VALUES (:name, :created_at, :updated_at)`

	row := categoryRow{
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	result, err := r.db.NamedExec(query, row)
	if err != nil {
		if isUniqueConstraint(err) {
			return domain.ErrNameDuplicate
		}
		return fmt.Errorf("category.Save: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("category.Save last id: %w", err)
	}

	c.ID = id
	return nil
}

func (r *categoryRepo) Update(c *domain.Category) error {
	query := `UPDATE categories
	          SET name = :name, updated_at = :updated_at
	          WHERE id = :id`

	row := categoryRow{
		ID:        c.ID,
		Name:      c.Name,
		UpdatedAt: c.UpdatedAt,
	}

	result, err := r.db.NamedExec(query, row)
	if err != nil {
		if isUniqueConstraint(err) {
			return domain.ErrNameDuplicate
		}
		return fmt.Errorf("category.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("category.Update rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *categoryRepo) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("category.Delete: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("category.Delete rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
