package gormdb_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"inventario/internal/platform/gormdb"

	"gorm.io/gorm"
)

// itemRepo is a minimal repository used to observe transaction effects. Like
// production repositories it resolves its handle from ctx.
type itemRepo struct{ db *gorm.DB }

func (r itemRepo) insert(ctx context.Context, name string) error {
	return gormdb.Conn(ctx, r.db).Exec(`INSERT INTO items (name) VALUES (?)`, name).Error
}

func newTestDB(t *testing.T) (*gorm.DB, *sql.DB) {
	t.Helper()
	// Open uses a single connection, which keeps :memory: shared.
	db, err := gormdb.Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql pool: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := db.Exec(`CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`).Error; err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db, sqlDB
}

func countItems(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Table("items").Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// withTimeout fails the test instead of hanging forever if fn deadlocks on
// the single connection.
func withTimeout(t *testing.T, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock: unit of work hung")
	}
}

func TestUnitOfWork_Do(t *testing.T) {
	ctx := context.Background()

	t.Run("commits when fn succeeds", func(t *testing.T) {
		db, _ := newTestDB(t)
		repo := itemRepo{db: db}

		err := gormdb.NewUnitOfWork(db).Do(ctx, func(ctx context.Context) error {
			if err := repo.insert(ctx, "a"); err != nil {
				return err
			}
			return repo.insert(ctx, "b")
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := countItems(t, db); got != 2 {
			t.Errorf("got %d items, want 2", got)
		}
	})

	t.Run("rolls back all writes when fn fails", func(t *testing.T) {
		db, _ := newTestDB(t)
		repo := itemRepo{db: db}
		fnErr := errors.New("second write failed")

		err := gormdb.NewUnitOfWork(db).Do(ctx, func(ctx context.Context) error {
			if err := repo.insert(ctx, "a"); err != nil {
				return err
			}
			return fnErr
		})

		if !errors.Is(err, fnErr) {
			t.Errorf("got %v, want %v", err, fnErr)
		}
		if got := countItems(t, db); got != 0 {
			t.Errorf("got %d items, want 0", got)
		}
	})

	t.Run("rolls back and re-panics when fn panics", func(t *testing.T) {
		db, _ := newTestDB(t)
		repo := itemRepo{db: db}

		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic to propagate")
				}
			}()
			_ = gormdb.NewUnitOfWork(db).Do(ctx, func(ctx context.Context) error {
				if err := repo.insert(ctx, "a"); err != nil {
					return err
				}
				panic("boom")
			})
		}()

		if got := countItems(t, db); got != 0 {
			t.Errorf("got %d items, want 0", got)
		}
	})

	t.Run("returns begin error without calling fn", func(t *testing.T) {
		db, sqlDB := newTestDB(t)
		_ = sqlDB.Close()
		called := false

		err := gormdb.NewUnitOfWork(db).Do(ctx, func(context.Context) error {
			called = true
			return nil
		})

		if err == nil {
			t.Fatal("expected begin error")
		}
		if called {
			t.Error("fn must not run when begin fails")
		}
	})

	t.Run("nested Do joins the outer transaction", func(t *testing.T) {
		db, _ := newTestDB(t)
		repo := itemRepo{db: db}
		u := gormdb.NewUnitOfWork(db)
		outerErr := errors.New("outer failed after inner succeeded")

		var err error
		withTimeout(t, func() {
			err = u.Do(ctx, func(ctx context.Context) error {
				if err := repo.insert(ctx, "outer"); err != nil {
					return err
				}
				if err := u.Do(ctx, func(ctx context.Context) error {
					return repo.insert(ctx, "inner")
				}); err != nil {
					return err
				}
				return outerErr
			})
		})

		if !errors.Is(err, outerErr) {
			t.Errorf("got %v, want %v", err, outerErr)
		}
		if got := countItems(t, db); got != 0 {
			t.Errorf("got %d items, want 0 (inner work must roll back with outer)", got)
		}
	})
}

func TestConn(t *testing.T) {
	db, _ := newTestDB(t)
	inTx := func(g *gorm.DB) bool {
		_, ok := g.Statement.ConnPool.(gorm.TxCommitter)
		return ok
	}

	if inTx(gormdb.Conn(context.Background(), db)) {
		t.Error("without a unit of work Conn must not return a transaction")
	}

	err := gormdb.NewUnitOfWork(db).Do(context.Background(), func(ctx context.Context) error {
		if !inTx(gormdb.Conn(ctx, db)) {
			t.Error("inside Do Conn must return the transaction")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
