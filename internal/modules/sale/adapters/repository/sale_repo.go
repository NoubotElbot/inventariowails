// Package repository is the GORM driven adapter implementing
// ports.SaleRepository for the sale module.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"

	"gorm.io/gorm"
)

type saleRepo struct {
	db *gorm.DB
}

var _ ports.SaleRepository = (*saleRepo)(nil)

func NewSaleRepository(db *gorm.DB) ports.SaleRepository {
	return &saleRepo{db: db}
}

// --- modelos de persistencia ---

type saleModel struct {
	ID        int64 `gorm:"primaryKey"`
	Total     int64
	Items     []saleItemModel `gorm:"foreignKey:SaleID"`
	CreatedAt time.Time
}

func (saleModel) TableName() string { return "sales" }

type saleItemModel struct {
	ID        int64 `gorm:"primaryKey"`
	SaleID    int64
	ProductID int64
	Product   *saleProductModel `gorm:"foreignKey:ProductID"`
	Quantity  int
	UnitPrice int64
	Subtotal  int64
	CreatedAt time.Time
}

func (saleItemModel) TableName() string { return "sale_items" }

// saleProductModel is a read-only projection of the products table: the
// belongs-to side of saleItemModel, loaded with Preload("Items.Product").
type saleProductModel struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
	Code string
}

func (saleProductModel) TableName() string { return "products" }

func (m *saleProductModel) toDomain() *domain.SaleProduct {
	if m == nil {
		return nil
	}
	return &domain.SaleProduct{ID: m.ID, Name: m.Name, Code: m.Code}
}

func (m saleModel) toDomain() domain.Sale {
	items := make([]domain.SaleItem, 0, len(m.Items))
	for _, it := range m.Items {
		items = append(items, domain.SaleItem{
			ID:        it.ID,
			SaleID:    it.SaleID,
			ProductID: it.ProductID,
			Product:   it.Product.toDomain(),
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			Subtotal:  it.Subtotal,
			CreatedAt: it.CreatedAt,
		})
	}
	return domain.Sale{
		ID:        m.ID,
		Total:     m.Total,
		Items:     items,
		CreatedAt: m.CreatedAt,
	}
}

// fromDomain builds a new row (IDs left to the database) for insertion.
func fromDomain(s *domain.Sale) saleModel {
	row := saleModel{
		Total:     s.Total,
		CreatedAt: s.CreatedAt,
		Items:     make([]saleItemModel, 0, len(s.Items)),
	}
	for _, it := range s.Items {
		row.Items = append(row.Items, saleItemModel{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			Subtotal:  it.Subtotal,
			CreatedAt: it.CreatedAt,
		})
	}
	return row
}

// --- implementación ---

// Save inserts the sale and its items in a single GORM Create, which is
// atomic on its own and joins the caller's transaction (via gormdb.Conn)
// when called inside a uow.Do.
func (r *saleRepo) Save(ctx context.Context, s *domain.Sale) error {
	row := fromDomain(s)

	if err := gormdb.Conn(ctx, r.db).Create(&row).Error; err != nil {
		return fmt.Errorf("sale.Save: %w", err)
	}

	s.ID = row.ID
	for i := range s.Items {
		s.Items[i].ID = row.Items[i].ID
		s.Items[i].SaleID = row.ID
	}
	return nil
}

func (r *saleRepo) FindAll(ctx context.Context, q domain.FindAllQuery) (pagination.Result[domain.Sale], error) {
	query := gormdb.Conn(ctx, r.db).
		Model(&saleModel{}).
		Preload("Items.Product").
		Order("created_at DESC, id DESC")

	res, err := pagination.Paginate[saleModel](query, q.Params)
	if err != nil {
		return pagination.Result[domain.Sale]{}, fmt.Errorf("sale.FindAll: %w", err)
	}

	sales := make([]domain.Sale, 0, len(res.Data))
	for _, row := range res.Data {
		sales = append(sales, row.toDomain())
	}
	return pagination.Result[domain.Sale]{Data: sales, Meta: res.Meta}, nil
}

func (r *saleRepo) FindByID(ctx context.Context, id int64) (*domain.Sale, error) {
	var row saleModel
	if err := gormdb.Conn(ctx, r.db).Preload("Items.Product").First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("sale.FindByID: %w", err)
	}

	s := row.toDomain()
	return &s, nil
}
