package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"inventario/database"
	"inventario/internal/modules/sale/adapters/repository"
	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"
)

// newTestRepo abre una SQLite en memoria con el esquema real (database.Migrate).
func newTestRepo(t *testing.T) (ports.SaleRepository, *gorm.DB) {
	t.Helper()
	// Open usa una sola conexión, lo que mantiene compartida la base :memory:.
	db, err := gormdb.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = gormdb.Close(db) })

	require.NoError(t, database.Migrate(db))
	return repository.NewSaleRepository(db), db
}

// seedProduct crea un producto para satisfacer la FK de sale_items.
func seedProduct(t *testing.T, db *gorm.DB, code string) int64 {
	t.Helper()
	row := database.Product{Name: "Producto " + code, Code: code, Price: 100, Stock: 10, Active: true}
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

// newSale arma una venta de dominio con un item por producto, cantidad 1 y
// precio unitario 100 (sin IDs, como la construye el servicio).
func newSale(createdAt time.Time, productIDs ...int64) *domain.Sale {
	s := &domain.Sale{CreatedAt: createdAt}
	for i, pid := range productIDs {
		qty := i + 1
		sub := int64(100 * qty)
		s.Items = append(s.Items, domain.SaleItem{
			ProductID: pid, Quantity: qty, UnitPrice: 100, Subtotal: sub, CreatedAt: createdAt,
		})
		s.Total += sub
	}
	return s
}

func TestSaleRepo_Save(t *testing.T) {
	t.Run("asigna IDs a la venta y a sus items y persiste los montos", func(t *testing.T) {
		repo, db := newTestRepo(t)
		p1 := seedProduct(t, db, "P1")
		p2 := seedProduct(t, db, "P2")
		now := time.Now()
		s := newSale(now, p1, p2)

		require.NoError(t, repo.Save(t.Context(), s))

		require.NotZero(t, s.ID)
		for _, it := range s.Items {
			assert.NotZero(t, it.ID)
			assert.Equal(t, s.ID, it.SaleID)
		}

		got, err := repo.FindByID(t.Context(), s.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(300), got.Total)
		assert.WithinDuration(t, now, got.CreatedAt, time.Second)
		require.Len(t, got.Items, 2)
		byProduct := map[int64]domain.SaleItem{}
		for _, it := range got.Items {
			byProduct[it.ProductID] = it
		}
		assert.Equal(t, 1, byProduct[p1].Quantity)
		assert.Equal(t, int64(100), byProduct[p1].UnitPrice)
		assert.Equal(t, int64(100), byProduct[p1].Subtotal)
		assert.Equal(t, 2, byProduct[p2].Quantity)
		assert.Equal(t, int64(200), byProduct[p2].Subtotal)
	})
}

func TestSaleRepo_FindByID(t *testing.T) {
	t.Run("carga el producto relacionado de cada item", func(t *testing.T) {
		repo, db := newTestRepo(t)
		p := seedProduct(t, db, "P1")
		s := newSale(time.Now(), p)
		require.NoError(t, repo.Save(t.Context(), s))

		got, err := repo.FindByID(t.Context(), s.ID)

		require.NoError(t, err)
		require.Len(t, got.Items, 1)
		assert.Equal(t, &domain.SaleProduct{ID: p, Name: "Producto P1", Code: "P1"}, got.Items[0].Product)
	})

	t.Run("ID inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		_, err := repo.FindByID(t.Context(), 999)

		require.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestSaleRepo_FindAll(t *testing.T) {
	t.Run("ordena por created_at DESC e incluye los items", func(t *testing.T) {
		repo, db := newTestRepo(t)
		p := seedProduct(t, db, "P1")
		base := time.Now().Add(-time.Hour)
		oldest := newSale(base, p)
		newest := newSale(base.Add(2*time.Minute), p, p)
		middle := newSale(base.Add(time.Minute), p)
		for _, s := range []*domain.Sale{oldest, newest, middle} {
			require.NoError(t, repo.Save(t.Context(), s))
		}

		got, err := repo.FindAll(t.Context(), pageQuery(1, 10))

		require.NoError(t, err)
		require.Len(t, got.Data, 3)
		assert.Equal(t, []int64{newest.ID, middle.ID, oldest.ID}, []int64{got.Data[0].ID, got.Data[1].ID, got.Data[2].ID})
		assert.Len(t, got.Data[0].Items, 2)
		assert.Len(t, got.Data[1].Items, 1)
	})

	t.Run("pagina los resultados", func(t *testing.T) {
		repo, db := newTestRepo(t)
		p := seedProduct(t, db, "P1")
		base := time.Now().Add(-time.Hour)
		ids := make([]int64, 0, 5)
		for i := range 5 {
			s := newSale(base.Add(time.Duration(i)*time.Minute), p)
			require.NoError(t, repo.Save(t.Context(), s))
			ids = append(ids, s.ID)
		}

		got, err := repo.FindAll(t.Context(), pageQuery(2, 2))

		require.NoError(t, err)
		require.Len(t, got.Data, 2)
		// Más nuevas primero: la página 2 trae la 3.ª y 2.ª ventas creadas.
		assert.Equal(t, []int64{ids[2], ids[1]}, []int64{got.Data[0].ID, got.Data[1].ID})
		assert.Equal(t, int64(5), got.Meta.TotalItems)
		assert.Equal(t, 3, got.Meta.TotalPages)
		assert.Equal(t, 2, got.Meta.Page)
	})

	t.Run("carga el producto relacionado de cada item", func(t *testing.T) {
		repo, db := newTestRepo(t)
		p := seedProduct(t, db, "P1")
		require.NoError(t, repo.Save(t.Context(), newSale(time.Now(), p)))

		got, err := repo.FindAll(t.Context(), pageQuery(1, 10))

		require.NoError(t, err)
		require.Len(t, got.Data, 1)
		require.Len(t, got.Data[0].Items, 1)
		assert.Equal(t, &domain.SaleProduct{ID: p, Name: "Producto P1", Code: "P1"}, got.Data[0].Items[0].Product)
	})

	t.Run("tabla vacía devuelve página vacía no nula", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		got, err := repo.FindAll(t.Context(), pageQuery(1, 10))

		require.NoError(t, err)
		assert.NotNil(t, got.Data)
		assert.Empty(t, got.Data)
		assert.Zero(t, got.Meta.TotalItems)
	})
}

func pageQuery(page, pageSize int) domain.FindAllQuery {
	return domain.FindAllQuery{Params: pagination.Params{Page: page, PageSize: pageSize}}
}

func TestSaleRepo_UsesTransactionFromContext(t *testing.T) {
	repo, db := newTestRepo(t)
	p := seedProduct(t, db, "P1")
	errBoom := errors.New("boom")

	err := gormdb.NewUnitOfWork(db).Do(t.Context(), func(ctx context.Context) error {
		require.NoError(t, repo.Save(ctx, newSale(time.Now(), p, p)))
		return errBoom
	})

	require.ErrorIs(t, err, errBoom)
	var sales, items int64
	require.NoError(t, db.Model(&database.Sale{}).Count(&sales).Error)
	require.NoError(t, db.Model(&database.SaleItem{}).Count(&items).Error)
	assert.Zero(t, sales, "la venta debe revertirse junto con la transacción")
	assert.Zero(t, items, "los items deben revertirse junto con la transacción")
}
