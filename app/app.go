package app

import (
	"fmt"
	"os"
	"path/filepath"

	"inventario/database"
	categoryhandler "inventario/internal/modules/category/adapters/handler"
	producthandler "inventario/internal/modules/product/adapters/handler"
	salehandler "inventario/internal/modules/sale/adapters/handler"
	stockhandler "inventario/internal/modules/stock/adapters/handler"
	systemhandler "inventario/internal/modules/system/adapters/handler"
	"inventario/internal/platform/gormdb"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB

	CategoryHandler *categoryhandler.CategoryHandler
	ProductHandler  *producthandler.ProductHandler
	StockHandler    *stockhandler.StockHandler
	SaleHandler     *salehandler.SaleHandler
	SystemHandler   *systemhandler.SystemHandler
}

func NewApp() (*App, error) {
	dbPath, err := ResolveDBPath()
	if err != nil {
		return nil, fmt.Errorf("resolviendo ruta de DB: %w", err)
	}

	db, err := gormdb.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("inicializando DB: %w", err)
	}

	if err := database.Migrate(db); err != nil {
		_ = gormdb.Close(db)
		return nil, fmt.Errorf("aplicando migraciones: %w", err)
	}

	a := &App{db: db}
	a.wire()

	return a, nil
}

func (a *App) Close() error {
	if a.db == nil {
		return nil
	}

	return gormdb.Close(a.db)
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
