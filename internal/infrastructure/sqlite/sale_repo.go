package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "inventario/internal/domain/sale"

	"github.com/jmoiron/sqlx"
)

type saleRepo struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) domain.Repository {
	return &saleRepo{db: db}
}

// --- modelos de persistencia ---

type saleRow struct {
	ID        int64     `db:"id"`
	Total     float64   `db:"total"`
	CreatedAt time.Time `db:"created_at"`
}

type saleItemRow struct {
	ID        int64     `db:"id"`
	SaleID    int64     `db:"sale_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	UnitPrice float64   `db:"unit_price"`
	Subtotal  float64   `db:"subtotal"`
	CreatedAt time.Time `db:"created_at"`
}

// --- implementación ---

func (r *saleRepo) Save(s *domain.Sale) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("sale.Save begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := saleRow{
		Total:     s.Total,
		CreatedAt: s.CreatedAt,
	}

	result, err := tx.NamedExec(
		`INSERT INTO sales (total, created_at) VALUES (:total, :created_at)`,
		row,
	)
	if err != nil {
		return fmt.Errorf("sale.Save insert sale: %w", err)
	}

	s.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("sale.Save last id: %w", err)
	}

	for i := range s.Items {
		s.Items[i].SaleID = s.ID
		itemRow := saleItemRow{
			SaleID:    s.ID,
			ProductID: s.Items[i].ProductID,
			Quantity:  s.Items[i].Quantity,
			UnitPrice: s.Items[i].UnitPrice,
			Subtotal:  s.Items[i].Subtotal,
			CreatedAt: s.Items[i].CreatedAt,
		}

		itemResult, err := tx.NamedExec(
			`INSERT INTO sale_items (sale_id, product_id, quantity, unit_price, subtotal, created_at)
			 VALUES (:sale_id, :product_id, :quantity, :unit_price, :subtotal, :created_at)`,
			itemRow,
		)
		if err != nil {
			return fmt.Errorf("sale.Save insert item: %w", err)
		}

		s.Items[i].ID, err = itemResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("sale.Save item last id: %w", err)
		}
	}

	return tx.Commit()
}

func (r *saleRepo) FindAll() ([]domain.Sale, error) {
	var rows []saleRow

	if err := r.db.Select(&rows,
		`SELECT id, total, created_at FROM sales ORDER BY created_at DESC`,
	); err != nil {
		return nil, fmt.Errorf("sale.FindAll: %w", err)
	}

	sales := make([]domain.Sale, 0, len(rows))
	for _, row := range rows {
		items, err := r.findItems(row.ID)
		if err != nil {
			return nil, err
		}
		sales = append(sales, domain.Sale{
			ID:        row.ID,
			Total:     row.Total,
			Items:     items,
			CreatedAt: row.CreatedAt,
		})
	}

	return sales, nil
}

func (r *saleRepo) FindByID(id int64) (*domain.Sale, error) {
	var row saleRow

	if err := r.db.Get(&row,
		`SELECT id, total, created_at FROM sales WHERE id = ?`, id,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("sale.FindByID: %w", err)
	}

	items, err := r.findItems(row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.Sale{
		ID:        row.ID,
		Total:     row.Total,
		Items:     items,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *saleRepo) findItems(saleID int64) ([]domain.SaleItem, error) {
	var rows []saleItemRow

	if err := r.db.Select(&rows,
		`SELECT id, sale_id, product_id, quantity, unit_price, subtotal, created_at
		 FROM sale_items WHERE sale_id = ?`, saleID,
	); err != nil {
		return nil, fmt.Errorf("sale.findItems: %w", err)
	}

	items := make([]domain.SaleItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.SaleItem{
			ID:        row.ID,
			SaleID:    row.SaleID,
			ProductID: row.ProductID,
			Quantity:  row.Quantity,
			UnitPrice: row.UnitPrice,
			Subtotal:  row.Subtotal,
			CreatedAt: row.CreatedAt,
		})
	}

	return items, nil
}
