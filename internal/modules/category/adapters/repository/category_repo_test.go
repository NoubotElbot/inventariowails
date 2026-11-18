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
	"inventario/internal/modules/category/adapters/repository"
	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
	"inventario/internal/platform/gormdb"
	"inventario/internal/platform/pagination"
)

// newTestRepo abre una SQLite en memoria con el esquema real (database.Migrate).
func newTestRepo(t *testing.T) (ports.CategoryRepository, *gorm.DB) {
	t.Helper()
	// Open usa una sola conexión, lo que mantiene compartida la base :memory:.
	db, err := gormdb.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = gormdb.Close(db) })

	require.NoError(t, database.Migrate(db))
	return repository.NewCategoryRepository(db), db
}

// seed inserta con el modelo del esquema (database.Category), no con el
// repositorio: así los tests detectan desfases entre schema.go y categoryModel.
func seed(t *testing.T, db *gorm.DB, name string, deleted bool) int64 {
	t.Helper()
	row := database.Category{Name: name}
	if deleted {
		now := time.Now()
		row.DeletedAt = &now
	}
	require.NoError(t, db.Create(&row).Error)
	return row.ID
}

func names(items []domain.Category) []string {
	out := make([]string, 0, len(items))
	for _, c := range items {
		out = append(out, c.Name)
	}
	return out
}

func TestCategoryRepo_FindAll(t *testing.T) {
	t.Run("filtra por estado y búsqueda", func(t *testing.T) {
		repo, db := newTestRepo(t)
		seed(t, db, "Electrónica", false)
		seed(t, db, "Ropa", false)
		seed(t, db, "Electro antiguo", true)

		tests := []struct {
			name  string
			query domain.FindAllQuery
			want  []string
		}{
			{"por defecto solo activas, id DESC", domain.FindAllQuery{}, []string{"Ropa", "Electrónica"}},
			{"OnlyInactive solo eliminadas", domain.FindAllQuery{OnlyInactive: true}, []string{"Electro antiguo"}},
			{"Search ignora mayúsculas", domain.FindAllQuery{Search: "elec"}, []string{"Electrónica"}},
			{"Search combinado con OnlyInactive", domain.FindAllQuery{Search: "elec", OnlyInactive: true}, []string{"Electro antiguo"}},
			{"Search sin coincidencias", domain.FindAllQuery{Search: "zzz"}, []string{}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := repo.FindAll(t.Context(), tt.query)

				require.NoError(t, err)
				assert.Equal(t, tt.want, names(got.Data))
				assert.Equal(t, int64(len(tt.want)), got.Meta.TotalItems)
			})
		}
	})

	t.Run("pagina y cuenta el total filtrado", func(t *testing.T) {
		repo, db := newTestRepo(t)
		for _, n := range []string{"A", "B", "C", "D", "E"} {
			seed(t, db, n, false)
		}
		seed(t, db, "Eliminada", true)

		tests := []struct {
			name   string
			params pagination.Params
			want   []string
		}{
			{"primera página", pagination.Params{Page: 1, PageSize: 2}, []string{"E", "D"}},
			{"última página incompleta", pagination.Params{Page: 3, PageSize: 2}, []string{"A"}},
			{"página fuera de rango", pagination.Params{Page: 4, PageSize: 2}, []string{}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := repo.FindAll(t.Context(), domain.FindAllQuery{Params: tt.params})

				require.NoError(t, err)
				assert.Equal(t, tt.want, names(got.Data))
				assert.Equal(t, int64(5), got.Meta.TotalItems)
				assert.Equal(t, 3, got.Meta.TotalPages)
				assert.Equal(t, tt.params.Page, got.Meta.Page)
				assert.Equal(t, tt.params.PageSize, got.Meta.PageSize)
			})
		}
	})
}

func TestCategoryRepo_FindOptions(t *testing.T) {
	t.Run("ordena por nombre e incluye eliminadas", func(t *testing.T) {
		repo, db := newTestRepo(t)
		ropa := seed(t, db, "Ropa", false)
		antigua := seed(t, db, "Antigua", true)
		elec := seed(t, db, "Electrónica", false)

		got, err := repo.FindOptions(t.Context())

		require.NoError(t, err)
		assert.Equal(t, []domain.Option{
			{ID: antigua, Name: "Antigua"},
			{ID: elec, Name: "Electrónica"},
			{ID: ropa, Name: "Ropa"},
		}, got)
	})

	t.Run("tabla vacía devuelve lista vacía", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		got, err := repo.FindOptions(t.Context())

		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestCategoryRepo_FindByID(t *testing.T) {
	t.Run("devuelve todos los campos, incluso si está eliminada", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, "Electrónica", true)

		got, err := repo.FindByID(t.Context(), id)

		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Electrónica", got.Name)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Minute)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Minute)
		require.NotNil(t, got.DeletedAt)
		assert.False(t, got.IsActive())
	})

	t.Run("ID inexistente devuelve ErrNotFound", func(t *testing.T) {
		repo, _ := newTestRepo(t)

		_, err := repo.FindByID(t.Context(), 999)

		require.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCategoryRepo_Save(t *testing.T) {
	t.Run("asigna ID y persiste los campos", func(t *testing.T) {
		repo, _ := newTestRepo(t)
		now := time.Now()
		c := &domain.Category{Name: "Electrónica", CreatedAt: now, UpdatedAt: now}

		require.NoError(t, repo.Save(t.Context(), c))
		require.NotZero(t, c.ID)

		got, err := repo.FindByID(t.Context(), c.ID)
		require.NoError(t, err)
		assert.Equal(t, "Electrónica", got.Name)
		assert.WithinDuration(t, now, got.CreatedAt, time.Second)
		assert.WithinDuration(t, now, got.UpdatedAt, time.Second)
		assert.Nil(t, got.DeletedAt)
	})

	t.Run("nombre duplicado devuelve ErrNameDuplicate", func(t *testing.T) {
		repo, db := newTestRepo(t)
		seed(t, db, "Ropa", false)

		err := repo.Save(t.Context(), &domain.Category{Name: "Ropa"})

		require.ErrorIs(t, err, domain.ErrNameDuplicate)
	})
}

func TestCategoryRepo_Update(t *testing.T) {
	t.Run("actualiza nombre y updated_at sin tocar created_at", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, "Viejo", false)
		before, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)

		c := *before
		c.Name = "Nuevo"
		c.UpdatedAt = before.UpdatedAt.Add(time.Hour)
		require.NoError(t, repo.Update(t.Context(), &c))

		after, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)
		assert.Equal(t, "Nuevo", after.Name)
		assert.True(t, after.UpdatedAt.After(before.UpdatedAt), "updated_at debe avanzar")
		assert.True(t, after.CreatedAt.Equal(before.CreatedAt), "created_at no debe cambiar")
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) *domain.Category
		wantErr error
	}{
		{
			name: "nombre duplicado devuelve ErrNameDuplicate",
			setup: func(t *testing.T, db *gorm.DB) *domain.Category {
				seed(t, db, "A", false)
				return &domain.Category{ID: seed(t, db, "B", false), Name: "A", UpdatedAt: time.Now()}
			},
			wantErr: domain.ErrNameDuplicate,
		},
		{
			name: "ID inexistente devuelve ErrNotFound",
			setup: func(*testing.T, *gorm.DB) *domain.Category {
				return &domain.Category{ID: 999, Name: "X", UpdatedAt: time.Now()}
			},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newTestRepo(t)
			c := tt.setup(t, db)

			require.ErrorIs(t, repo.Update(t.Context(), c), tt.wantErr)
		})
	}
}

func TestCategoryRepo_SoftDelete(t *testing.T) {
	t.Run("marca deleted_at", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, "Electrónica", false)

		require.NoError(t, repo.SoftDelete(t.Context(), id))

		got, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)
		assert.NotNil(t, got.DeletedAt)
	})

	tests := []struct {
		name  string
		setup func(t *testing.T, db *gorm.DB) int64
	}{
		{"ID inexistente", func(*testing.T, *gorm.DB) int64 { return 999 }},
		{"ya eliminada", func(t *testing.T, db *gorm.DB) int64 { return seed(t, db, "X", true) }},
	}
	for _, tt := range tests {
		t.Run(tt.name+" devuelve ErrNotFound", func(t *testing.T) {
			repo, db := newTestRepo(t)
			id := tt.setup(t, db)

			require.ErrorIs(t, repo.SoftDelete(t.Context(), id), domain.ErrNotFound)
		})
	}
}

func TestCategoryRepo_Restore(t *testing.T) {
	t.Run("limpia deleted_at", func(t *testing.T) {
		repo, db := newTestRepo(t)
		id := seed(t, db, "Electrónica", true)

		require.NoError(t, repo.Restore(t.Context(), id))

		got, err := repo.FindByID(t.Context(), id)
		require.NoError(t, err)
		assert.Nil(t, got.DeletedAt)
	})

	tests := []struct {
		name  string
		setup func(t *testing.T, db *gorm.DB) int64
	}{
		{"ID inexistente", func(*testing.T, *gorm.DB) int64 { return 999 }},
		{"no eliminada", func(t *testing.T, db *gorm.DB) int64 { return seed(t, db, "X", false) }},
	}
	for _, tt := range tests {
		t.Run(tt.name+" devuelve ErrNotFound", func(t *testing.T) {
			repo, db := newTestRepo(t)
			id := tt.setup(t, db)

			require.ErrorIs(t, repo.Restore(t.Context(), id), domain.ErrNotFound)
		})
	}
}

func TestCategoryRepo_UsesTransactionFromContext(t *testing.T) {
	repo, db := newTestRepo(t)
	errBoom := errors.New("boom")

	err := gormdb.NewUnitOfWork(db).Do(t.Context(), func(ctx context.Context) error {
		require.NoError(t, repo.Save(ctx, &domain.Category{Name: "Electrónica"}))
		return errBoom
	})

	require.ErrorIs(t, err, errBoom)
	var count int64
	require.NoError(t, db.Model(&database.Category{}).Count(&count).Error)
	assert.Zero(t, count, "el Save debe revertirse junto con la transacción")
}
