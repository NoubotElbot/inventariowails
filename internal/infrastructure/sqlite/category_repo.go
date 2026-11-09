package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "inventario/internal/domain/category"
	"inventario/internal/domain/pagination"

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
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (r categoryRow) toDomain() *domain.Category {
	return &domain.Category{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}

// --- implementación ---

func (r *categoryRepo) FindAll(q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	p := q.Params.Normalize()

	var wb whereBuilder
	if q.OnlyInactive {
		wb.And("deleted_at IS NOT NULL")
	} else {
		wb.And("deleted_at IS NULL")
	}
	if q.Search != "" {
		wb.And("name LIKE ?", "%"+q.Search+"%")
	}

	var total int64
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM categories "+wb.SQL(), wb.Args()...); err != nil {
		return pagination.Result[domain.Category]{}, fmt.Errorf("category.FindAll count: %w", err)
	}

	const cols = "id, name, created_at, updated_at, deleted_at"
	var rows []categoryRow
	if err := r.db.Select(&rows,
		"SELECT "+cols+" FROM categories "+wb.SQL()+" ORDER BY id DESC LIMIT ? OFFSET ?",
		append(wb.Args(), p.Limit(), p.Offset())...); err != nil {
		return pagination.Result[domain.Category]{}, fmt.Errorf("category.FindAll: %w", err)
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, *row.toDomain())
	}

	return pagination.NewResult(categories, total, p), nil
}

func (r *categoryRepo) FindOptions() ([]domain.Option, error) {
	var rows []categoryRow
	if err := r.db.Select(&rows, `SELECT id, name FROM categories ORDER BY name ASC`); err != nil {
		return nil, fmt.Errorf("category.FindOptions: %w", err)
	}

	options := make([]domain.Option, 0, len(rows))
	for _, row := range rows {
		options = append(options, domain.Option{ID: row.ID, Name: row.Name})
	}
	return options, nil
}

func (r *categoryRepo) FindByID(id int64) (*domain.Category, error) {
	var row categoryRow

	query := `SELECT id, name, created_at, updated_at, deleted_at
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

func (r *categoryRepo) SoftDelete(id int64) error {
	query := `UPDATE categories SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("category.SoftDelete: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("category.SoftDelete rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *categoryRepo) Restore(id int64) error {
	query := `UPDATE categories SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("category.Restore: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("category.Restore rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
