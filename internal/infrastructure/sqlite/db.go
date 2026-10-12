package sqlite

import (
	"embed"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var migrations embed.FS

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dsn+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite no soporta escrituras concurrentes

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	log.Println("SQLite conectado y migraciones aplicadas")
	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(db.DB, "database/migrations")
}
