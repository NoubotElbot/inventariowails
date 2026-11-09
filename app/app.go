package app

import (
	"fmt"
	"os"
	"path/filepath"

	"inventario/database"
	categoryhandler "inventario/internal/modules/category/handler"
	producthandler "inventario/internal/modules/product/handler"
	salehandler "inventario/internal/modules/sale/handler"
	stockhandler "inventario/internal/modules/stock/handler"
	"inventario/internal/platform/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type App struct {
	db *sqlx.DB

	CategoryHandler *categoryhandler.CategoryHandler
	ProductHandler  *producthandler.ProductHandler
	StockHandler    *stockhandler.StockHandler
	SaleHandler     *salehandler.SaleHandler
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
	_ = godotenv.Load()

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
