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
	"inventario/internal/modules/product/adapters/repository"
	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"
)

// newTestRepo abre una SQLite en memoria con el esquema real (database.Migrate).
func newTestRepo(t *testing.T) (ports.ProductRepository, *gorm.DB) {
	t.Helper()
	// Open usa una sola conexión, lo que mantiene compartida la base :memory:.
	db, err := gormdb.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = gormdb.Close(db) })

	require.NoError(t, database.Migrate(db))
	return repository.NewProductRepository(db), db
}

// seedCategory crea una categoría para satisfacer la FK de products.
func seedCategory(t *testing.T, db *gorm.DB, name string) int64 {
	t.Helper()
	row := database.Category{Name: name}
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

// seed inserta con el modelo del esquema (database.Product), no con el
// repositorio: así los tests detectan desfases entre schema.go y productModel.
func seed(t *testing.T, db *gorm.DB, row database.Product) int64 {
	t.Helper()
	if row.Price == 0 {
		row.Price = 100
	}
	row.Active = true
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

func ptr[T any](v T) *T { return &v }

func codes(items []domain.Product) []string {
	out := make([]string, 0, len(items))
	for _, p := range items {
		out = append(out, p.Code)
	}
	return out
}

func TestProductRepo_FindAll(t *testing.T) {
	t.Run("filtra por búsqueda y categorías, ordena por nombre", func(t *testing.T) {
		repo, db := newTestRepo(t)
		catA := seedCategory(t, db, "A")
		catB := seedCategory(t, db, "B")
		catC := seedCategory(t, db, "C")
		seed(t, db, database.Product{Name: "Mouse USB", Code: "MOU-001", Barcode: ptr("789000"), CategoryID: &catA})
		seed(t, db, database.Product{Name: "Laptop Gamer", Code: "LAP-001", Barcode: ptr("123456"), CategoryID: &catA})
		seed(t, db, database.Product{Name: "Laptop Básica", Code: "LAP-002", CategoryID: &catB})
		seed(t, db, database.Product{Name: "Escritorio", Code: "DSK-001", CategoryID: &catC})
		seed(t, db, database.Product{Name: "Sin categoría", Code: "NOC-001"})

		tests := []struct {
			name  string
			query domain.FindAllQuery
			want  []string
		}{
			{"sin filtros, name ASC", domain.FindAllQuery{}, []string{"DSK-001", "LAP-002", "LAP-001", "MOU-001", "NOC-001"}},
			{"search por nombre ignora mayúsculas", domain.FindAllQuery{Search: "laptop"}, []string{"LAP-002", "LAP-001"}},
			{"search por código", domain.FindAllQuery{Search: "MOU"}, []string{"MOU-001"}},
			{"search por barcode", domain.FindAllQuery{Search: "1234"}, []string{"LAP-001"}},
			{"una categoría", domain.FindAllQuery{CategoryIDs: []int64{catA}}, []string{"LAP-001", "MOU-001"}},
			{"varias categorías", domain.FindAllQuery{CategoryIDs: []int64{catB, catC}}, []string{"DSK-001", "LAP-002"}},
			{"search combinado con categoría", domain.FindAllQuery{Search: "laptop", CategoryIDs: []int64{catA}}, []string{"LAP-001"}},
			{"search sin coincidencias", domain.FindAllQuery{Search: "zzz"}, []string{}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := repo.FindAll(t.Context(), tt.query)

				require.NoError(t, err)
				assert.Equal(t, tt.want, codes(got.Data))
				assert.Equal(t, int64(len(tt.want)), got.Meta.TotalItems)
			})
		}
	})

	t.Run("mapea todos los campos", func(t *testing.T) {
		repo, db := newTestRepo(t)
		cat := seedCategory(t, db, "A")
		id := seed(t, db, database.Product{
			Name: "Laptop", Code: "LAP-001", Barcode: ptr("123"), CategoryID: &cat,
			Price: 999, Stock: 4, MinStock: 5,
		})

		got, err := repo.FindAll(t.Context(), domain.FindAllQuery{})

		require.NoError(t, err)
		require.Len(t, got.Data, 1)
		p := got.Data[0]
		assert.Equal(t, id, p.ID)
		assert.Equal(t, &cat, p.CategoryID)
		assert.Equal(t, "123", p.Barcode)
		assert.Equal(t, int64(999), p.Price)
		assert.Equal(t, 4, p.Stock)
		assert.Equal(t, 5, p.MinStock)
		assert.True(t, p.Active)
	})

	t.Run("pagina y cuenta el total filtrado", func(t *testing.T) {
		repo, db := newTestRepo(t)
		for _, n := range []string{"A", "B", "C", "D", "E"} {
			seed(t, db, database.Product{Name: n, Code: n})
		}

		tests := []struct {
			name   string
			params pagination.Params
			want   []string
		}{
			{"primera página", pagination.Params{Page: 1, PageSize: 2}, []string{"A", "B"}},
			{"última página incompleta", pagination.Params{Page: 3, PageSize: 2}, []string{"E"}},
			{"página fuera de rango", pagination.Params{Page: 4, PageSize: 2}, []string{}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := repo.FindAll(t.Context(), domain.FindAllQuery{Params: tt.params})

				require.NoError(t, err)
				assert.Equal(t, tt.want, codes(got.Data))
				assert.Equal(t, int64(5), got.Meta.TotalItems)
				assert.Equal(t, 3, got.Meta.TotalPages)
				assert.Equal(t, tt.params.Page, got.Meta.Page)
				assert.Equal(t, tt.params.PageSize, got.Meta.PageSize)
			})
		}
	})
}

func TestProductRepo_FindByID(t *testing.T) {
	t.Run("devuelve todos los campos", func(t *testing.T) {
		repo, db := newTestRepo(t)
		cat := seedCategory(t, db, "A")
		id := seed(t, db, database.Product{
			Name: "Laptop", Code: "LAP-001", Barcode: ptr("123"), CategoryID: &cat,
			Price: 999, Stock: 4, MinStock: 5,
		})

		got, err := repo.FindByID(t.Context(), id)

		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, &cat, got.CategoryID)
		assert.Equal(t, "Laptop", got.Name)
		assert.Equal(t, "LAP-001", got.Code)
		assert.Equal(t, "123", got.Barcode)
		assert.Equal(t, int64(999), got.Price)
		assert.Equal(t, 4, got.Stock)
		assert.Equal(t, 5, got.MinStock)
		assert.True(t, got.Active)
		assert.True(t, got.HasLowStock())
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Minute)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Minute)
	})

	t.Run("barcode NULL se lee como cadena vacía", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001"})

		got, err := repo.FindByID(t.Context(), id)

		require.NoError(t, err)
		assert.Empty(t, got.Barcode)
		assert.Nil(t, got.CategoryID)
	})

	t.Run("ID inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		_, err := repo.FindByID(t.Context(), 999)

		require.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestProductRepo_FindByCode(t *testing.T) {
	t.Run("encuentra por código exacto", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001"})

		got, err := repo.FindByCode(t.Context(), "LAP-001")

		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("código inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, db := newTestRepo(t)
		seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001"})

		_, err := repo.FindByCode(t.Context(), "LAP")

		require.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestProductRepo_Save(t *testing.T) {
	t.Run("asigna ID y persiste los campos", func(t *testing.T) {
		repo, db := newTestRepo(t)
		cat := seedCategory(t, db, "A")
		now := time.Now()
		p := &domain.Product{
			CategoryID: &cat, Name: "Laptop", Code: "LAP-001", Barcode: "123",
			Price: 999, Stock: 10, MinStock: 2, Active: true, CreatedAt: now, UpdatedAt: now,
		}

		require.NoError(t, repo.Save(t.Context(), p))
		require.NotZero(t, p.ID)

		got, err := repo.FindByID(t.Context(), p.ID)
		require.NoError(t, err)
		assert.Equal(t, &cat, got.CategoryID)
		assert.Equal(t, "Laptop", got.Name)
		assert.Equal(t, "123", got.Barcode)
		assert.Equal(t, int64(999), got.Price)
		assert.Equal(t, 10, got.Stock)
		assert.Equal(t, 2, got.MinStock)
		assert.True(t, got.Active)
		assert.WithinDuration(t, now, got.CreatedAt, time.Second)
	})

	t.Run("código duplicado devuelve ErrCodeDuplicate", func(t *testing.T) {
		repo, db := newTestRepo(t)
		seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001"})

		err := repo.Save(t.Context(), &domain.Product{Name: "Otra", Code: "LAP-001", Price: 1})

		require.ErrorIs(t, err, domain.ErrCodeDuplicate)
	})
}

func TestProductRepo_Update(t *testing.T) {
	t.Run("actualiza campos editables sin tocar stock ni created_at", func(t *testing.T) {
		repo, db := newTestRepo(t)
		cat := seedCategory(t, db, "A")
		id := seed(t, db, database.Product{Name: "Viejo", Code: "OLD", Stock: 10})
		before, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)

		p := *before
		p.CategoryID = &cat
		p.Name = "Nuevo"
		p.Code = "NEW"
		p.Barcode = "999"
		p.Price = 500
		p.MinStock = 3
		p.Active = false
		p.Stock = 99 // no debe persistirse
		p.CreatedAt = before.CreatedAt.Add(-time.Hour)
		p.UpdatedAt = before.UpdatedAt.Add(time.Hour)
		require.NoError(t, repo.Update(t.Context(), &p))

		after, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)
		assert.Equal(t, &cat, after.CategoryID)
		assert.Equal(t, "Nuevo", after.Name)
		assert.Equal(t, "NEW", after.Code)
		assert.Equal(t, "999", after.Barcode)
		assert.Equal(t, int64(500), after.Price)
		assert.Equal(t, 3, after.MinStock)
		assert.False(t, after.Active, "active=false debe persistirse")
		assert.Equal(t, 10, after.Stock, "stock solo cambia vía UpdateStock")
		assert.True(t, after.CreatedAt.Equal(before.CreatedAt), "created_at no debe cambiar")
		assert.True(t, after.UpdatedAt.After(before.UpdatedAt), "updated_at debe avanzar")
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) *domain.Product
		wantErr error
	}{
		{
			name: "código duplicado devuelve ErrCodeDuplicate",
			setup: func(t *testing.T, db *gorm.DB) *domain.Product {
				seed(t, db, database.Product{Name: "A", Code: "A"})
				id := seed(t, db, database.Product{Name: "B", Code: "B"})
				return &domain.Product{ID: id, Name: "B", Code: "A", Price: 1, UpdatedAt: time.Now()}
			},
			wantErr: domain.ErrCodeDuplicate,
		},
		{
			name: "ID inexistente devuelve ErrNotFound",
			setup: func(*testing.T, *gorm.DB) *domain.Product {
				return &domain.Product{ID: 999, Name: "X", Code: "X", Price: 1, UpdatedAt: time.Now()}
			},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newTestRepo(t)
			p := tt.setup(t, db)

			require.ErrorIs(t, repo.Update(t.Context(), p), tt.wantErr)
		})
	}
}

func TestProductRepo_Delete(t *testing.T) {
	t.Run("elimina la fila", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001"})

		require.NoError(t, repo.Delete(t.Context(), id))

		_, err := repo.FindByID(t.Context(), id)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("ID inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		require.ErrorIs(t, repo.Delete(t.Context(), 999), domain.ErrNotFound)
	})
}

func TestProductRepo_UpdateStock(t *testing.T) {
	t.Run("actualiza solo el stock", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, database.Product{Name: "Laptop", Code: "LAP-001", Stock: 10})

		require.NoError(t, repo.UpdateStock(t.Context(), id, 3))

		got, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)
		assert.Equal(t, 3, got.Stock)
		assert.Equal(t, "Laptop", got.Name)
	})

	t.Run("ID inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		require.ErrorIs(t, repo.UpdateStock(t.Context(), 999, 1), domain.ErrNotFound)
	})
}

func TestProductRepo_UsesTransactionFromContext(t *testing.T) {
	repo, db := newTestRepo(t)
	errBoom := errors.New("boom")

	err := gormdb.NewUnitOfWork(db).Do(t.Context(), func(ctx context.Context) error {
		require.NoError(t, repo.Save(ctx, &domain.Product{Name: "Laptop", Code: "LAP-001", Price: 1, Active: true}))
		return errBoom
	})

	require.ErrorIs(t, err, errBoom)
	var count int64
	require.NoError(t, db.Model(&database.Product{}).Count(&count).Error)
	assert.Zero(t, count, "el Save debe revertirse junto con la transacción")
}
