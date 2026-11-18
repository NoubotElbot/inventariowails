// Package repository is the GORM driven adapter implementing
// ports.StockRepository for the stock module.
package repository

import (
	"context"
	"fmt"
	"time"

	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/gormdb"

	"gorm.io/gorm"
)

type stockRepo struct {
	db *gorm.DB
}

var _ ports.StockRepository = (*stockRepo)(nil)

// NewStockRepository returns the GORM-backed stock repository. Each call
// resolves its handle with gormdb.Conn, so it joins the transaction of a
// surrounding uow.Do when there is one.
func NewStockRepository(db *gorm.DB) ports.StockRepository {
	return &stockRepo{db: db}
}

// --- modelo de persistencia ---

// stockMovementModel maps the stock_movements table.
type stockMovementModel struct {
	ID           int64 `gorm:"primaryKey"`
	ProductID    int64
	Quantity     int
	MovementType string
	Notes        string
	CreatedAt    time.Time
}

func (stockMovementModel) TableName() string { return "stock_movements" }

func (m stockMovementModel) toDomain() domain.Movement {
	return domain.Movement{
		ID:           m.ID,
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		MovementType: domain.MovementType(m.MovementType),
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
	}
}

func fromDomain(m *domain.Movement) stockMovementModel {
	return stockMovementModel{
		ID:           m.ID,
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		MovementType: string(m.MovementType),
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
	}
}

// --- implementación ---

func (r *stockRepo) Save(ctx context.Context, m *domain.Movement) error {
	row := fromDomain(m)
	row.ID = 0

	if err := gormdb.Conn(ctx, r.db).Create(&row).Error; err != nil {
		return fmt.Errorf("stock.Save: %w", err)
	}

	m.ID = row.ID
	return nil
}

func (r *stockRepo) FindByProductID(ctx context.Context, productID int64) ([]domain.Movement, error) {
	var rows []stockMovementModel

	err := gormdb.Conn(ctx, r.db).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("stock.FindByProductID: %w", err)
	}

	movements := make([]domain.Movement, 0, len(rows))
	for _, row := range rows {
		movements = append(movements, row.toDomain())
	}
	return movements, nil
}
