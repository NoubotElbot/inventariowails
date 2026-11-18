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
	"inventario/internal/modules/stock/adapters/repository"
	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/gormdb"
)

// newTestRepo abre una SQLite en memoria con el esquema real (database.Migrate).
func newTestRepo(t *testing.T) (ports.StockRepository, *gorm.DB) {
	t.Helper()
	// Open usa una sola conexión, lo que mantiene compartida la base :memory:.
	db, err := gormdb.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = gormdb.Close(db) })

	require.NoError(t, database.Migrate(db))
	return repository.NewStockRepository(db), db
}

// seedProduct crea un producto para satisfacer la FK de stock_movements.
func seedProduct(t *testing.T, db *gorm.DB, code string) int64 {
	t.Helper()
	row := database.Product{Name: "Producto " + code, Code: code, Price: 100, Active: true}
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

// seedMovement inserta con el modelo del esquema (database.StockMovement), no
// con el repositorio: así los tests detectan desfases entre schema.go y
// stockMovementModel.
func seedMovement(t *testing.T, db *gorm.DB, row database.StockMovement) int64 {
	t.Helper()
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

func TestStockRepo_Save(t *testing.T) {
	t.Run("asigna ID y persiste los campos", func(t *testing.T) {
		repo, db := newTestRepo(t)
		pid := seedProduct(t, db, "P1")
		now := time.Now()
		m := &domain.Movement{
			ProductID: pid, Quantity: 3, MovementType: domain.MovementAdjustmentIn,
			Notes: "conteo", CreatedAt: now,
		}

		require.NoError(t, repo.Save(t.Context(), m))
		require.NotZero(t, m.ID)

		got, err := repo.FindByProductID(t.Context(), pid)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, m.ID, got[0].ID)
		assert.Equal(t, pid, got[0].ProductID)
		assert.Equal(t, 3, got[0].Quantity)
		assert.Equal(t, domain.MovementAdjustmentIn, got[0].MovementType)
		assert.Equal(t, "conteo", got[0].Notes)
		assert.WithinDuration(t, now, got[0].CreatedAt, time.Second)
	})

	t.Run("ignora un ID previo y deja que la base lo asigne", func(t *testing.T) {
		repo, db := newTestRepo(t)
		pid := seedProduct(t, db, "P1")
		first := &domain.Movement{ProductID: pid, Quantity: 1, MovementType: domain.MovementSale, CreatedAt: time.Now()}
		require.NoError(t, repo.Save(t.Context(), first))

		second := &domain.Movement{ID: first.ID, ProductID: pid, Quantity: 1, MovementType: domain.MovementSale, CreatedAt: time.Now()}
		require.NoError(t, repo.Save(t.Context(), second))

		assert.NotEqual(t, first.ID, second.ID)
	})
}

func TestStockRepo_FindByProductID(t *testing.T) {
	t.Run("filtra por producto y ordena por created_at DESC", func(t *testing.T) {
		repo, db := newTestRepo(t)
		pid := seedProduct(t, db, "P1")
		other := seedProduct(t, db, "P2")
		base := time.Now().Add(-time.Hour)
		oldest := seedMovement(t, db, database.StockMovement{ProductID: pid, Quantity: 1, MovementType: "sale", CreatedAt: base})
		newest := seedMovement(t, db, database.StockMovement{ProductID: pid, Quantity: 2, MovementType: "return", CreatedAt: base.Add(2 * time.Minute)})
		middle := seedMovement(t, db, database.StockMovement{ProductID: pid, Quantity: 3, MovementType: "adjustment_in", CreatedAt: base.Add(time.Minute)})
		seedMovement(t, db, database.StockMovement{ProductID: other, Quantity: 9, MovementType: "sale", CreatedAt: base})

		got, err := repo.FindByProductID(t.Context(), pid)

		require.NoError(t, err)
		ids := make([]int64, 0, len(got))
		for _, m := range got {
			ids = append(ids, m.ID)
		}
		assert.Equal(t, []int64{newest, middle, oldest}, ids)
	})

	t.Run("notes NULL se lee como cadena vacía", func(t *testing.T) {
		repo, db := newTestRepo(t)
		pid := seedProduct(t, db, "P1")
		seedMovement(t, db, database.StockMovement{ProductID: pid, Quantity: 1, MovementType: "sale", CreatedAt: time.Now()})

		got, err := repo.FindByProductID(t.Context(), pid)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Empty(t, got[0].Notes)
	})

	t.Run("sin movimientos devuelve lista vacía no nula", func(t *testing.T) {
		repo, db := newTestRepo(t)
		pid := seedProduct(t, db, "P1")

		got, err := repo.FindByProductID(t.Context(), pid)

		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})
}

func TestStockRepo_UsesTransactionFromContext(t *testing.T) {
	repo, db := newTestRepo(t)
	pid := seedProduct(t, db, "P1")
	errBoom := errors.New("boom")

	err := gormdb.NewUnitOfWork(db).Do(t.Context(), func(ctx context.Context) error {
		require.NoError(t, repo.Save(ctx, &domain.Movement{
			ProductID: pid, Quantity: 1, MovementType: domain.MovementSale, CreatedAt: time.Now(),
		}))
		return errBoom
	})

	require.ErrorIs(t, err, errBoom)
	var count int64
	require.NoError(t, db.Model(&database.StockMovement{}).Count(&count).Error)
	assert.Zero(t, count, "el Save debe revertirse junto con la transacción")
}
