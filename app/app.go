package app

import (
	"fmt"
	"os"
	"path/filepath"

	"inventario/database"
	"inventario/internal/handler"
	"inventario/internal/infrastructure/sqlite"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db *sqlx.DB

	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	StockHandler    *handler.StockHandler
	SaleHandler     *handler.SaleHandler
}

func NewApp() (*App, error) {
	dbPath, err := ResolveDBPath()
	if err != nil {
		return nil, fmt.Errorf("resolviendo ruta de DB: %w", err)
	}

	db, err := sqlite.NewDB(dbPath, database.Migrations)
	if err != nil {
		return nil, fmt.Errorf("inicializando DB: %w", err)
	}

	a := &App{db: db}
	a.wire()

	return a, nil
}

func (a *App) Close() error {
	if a.db == nil {
		return nil
	}

	return a.db.Close()
}

func ResolveDBPath() (string, error) {
	if p := os.Getenv("DB_PATH"); p != "" {
		return p, nil
	}
	if os.Getenv("ENV") == "development" {
		return "inventario_dev.db", nil
	}

	dataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dataDir, "inventario")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "inventario.db"), nil
}
