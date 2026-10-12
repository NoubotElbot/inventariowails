package sqlite

import (
	"fmt"
	"time"

	domain "inventario/internal/domain/stock"

	"github.com/jmoiron/sqlx"
)

type stockRepo struct {
	db *sqlx.DB
}

func NewStockRepository(db *sqlx.DB) domain.Repository {
	return &stockRepo{db: db}
}

type stockMovementRow struct {
	ID           int64     `db:"id"`
	ProductID    int64     `db:"product_id"`
	Quantity     int       `db:"quantity"`
	MovementType string    `db:"movement_type"`
	Notes        string    `db:"notes"`
	CreatedAt    time.Time `db:"created_at"`
}

func (r stockMovementRow) toDomain() *domain.Movement {
	return &domain.Movement{
		ID:           r.ID,
		ProductID:    r.ProductID,
		Quantity:     r.Quantity,
		MovementType: domain.MovementType(r.MovementType),
		Notes:        r.Notes,
		CreatedAt:    r.CreatedAt,
	}
}

func (r *stockRepo) Save(m *domain.Movement) error {
	query := `INSERT INTO stock_movements (product_id, quantity, movement_type, notes, created_at)
	          VALUES (:product_id, :quantity, :movement_type, :notes, :created_at)`

	row := stockMovementRow{
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		MovementType: string(m.MovementType),
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
	}

	result, err := r.db.NamedExec(query, row)
	if err != nil {
		return fmt.Errorf("stock.Save: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("stock.Save last id: %w", err)
	}

	m.ID = id
	return nil
}

func (r *stockRepo) FindByProductID(productID int64) ([]domain.Movement, error) {
	var rows []stockMovementRow

	query := `SELECT id, product_id, quantity, movement_type, notes, created_at
	          FROM stock_movements
	          WHERE product_id = ?
	          ORDER BY created_at DESC`

	if err := r.db.Select(&rows, query, productID); err != nil {
		return nil, fmt.Errorf("stock.FindByProductID: %w", err)
	}

	movements := make([]domain.Movement, 0, len(rows))
	for _, row := range rows {
		movements = append(movements, *row.toDomain())
	}

	return movements, nil
}
