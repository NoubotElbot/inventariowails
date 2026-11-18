// Package gormdb owns the app's SQLite connection, exposed as a GORM handle.
package gormdb

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Pure-Go SQLite driver, registered as "sqlite".
	_ "modernc.org/sqlite"
)

// Open opens the SQLite database at path through GORM using the pure-Go
// modernc driver. Schema is applied separately (see database.Migrate).
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        path + "?_foreign_keys=on&_journal_mode=WAL",
	}), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Warn,
			// Repositories turn ErrRecordNotFound into domain ErrNotFound;
			// it is an expected outcome (e.g. code-uniqueness checks), not
			// something to log.
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql pool: %w", err)
	}

	// SQLite allows a single writer; one connection also keeps :memory:
	// databases shared across queries.
	sqlDB.SetMaxOpenConns(1)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return db, nil
}

// Close releases the connection pool behind db.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
