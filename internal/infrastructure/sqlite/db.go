package sqlite

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func NewDB(dsn string, migrations fs.FS) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dsn+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := runMigrations(db, migrations); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	log.Println("SQLite conectado y migraciones aplicadas")
	return db, nil
}

func runMigrations(db *sqlx.DB, migrations fs.FS) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(db.DB, "migrations")
}
