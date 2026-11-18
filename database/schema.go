// Package database defines the persistent schema as GORM models and applies
// it with AutoMigrate.
//
// These models are the source of truth for the tables, constraints and
// foreign keys. Repositories keep their own (partial) models mapped to the
// same tables, so modules stay decoupled from each other.
//
// AutoMigrate only adds what is missing (tables, columns, indexes,
// constraints); it never drops or renames columns nor migrates data. Such
// changes must be added to Migrate as explicit steps using db.Migrator().
package database

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	Name      string     `gorm:"not null;unique"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time // soft delete handled by the category repository
}

func (Category) TableName() string { return "categories" }

type Product struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	CategoryID *int64    `gorm:"index"`
	Category   *Category `gorm:"constraint:OnDelete:SET NULL"`
	Name       string    `gorm:"not null"`
	Code       string    `gorm:"not null;unique"`
	Barcode    *string
	Price      int64     `gorm:"not null;default:0"`
	Stock      int       `gorm:"not null;default:0"`
	MinStock   int       `gorm:"not null;default:0"`
	Active     bool      `gorm:"type:integer;not null;default:1"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Product) TableName() string { return "products" }

type Sale struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	Total     int64      `gorm:"not null;check:total >= 0"`
	Items     []SaleItem `gorm:"foreignKey:SaleID"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Sale) TableName() string { return "sales" }

type SaleItem struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	SaleID    int64     `gorm:"not null;index"`
	ProductID int64     `gorm:"not null;index"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Quantity  int       `gorm:"not null;check:quantity > 0"`
	UnitPrice int64     `gorm:"not null;check:unit_price >= 0"`
	Subtotal  int64     `gorm:"not null;check:subtotal >= 0"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (SaleItem) TableName() string { return "sale_items" }

type StockMovement struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	ProductID    int64     `gorm:"not null;index"`
	Product      Product   `gorm:"foreignKey:ProductID"`
	Quantity     int       `gorm:"not null;check:quantity > 0"`
	Notes        *string
	MovementType string    `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (StockMovement) TableName() string { return "stock_movements" }

// Migrate creates or updates the schema. GORM orders the tables by their
// foreign-key dependencies.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Category{},
		&Product{},
		&Sale{},
		&SaleItem{},
		&StockMovement{},
	)
}
