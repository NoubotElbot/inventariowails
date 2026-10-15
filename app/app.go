package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"inventario/database"
	"inventario/internal/handler"
	"inventario/internal/infrastructure/sqlite"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db              *sqlx.DB
	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	StockHandler    *handler.StockHandler
}

func NewApp() *App {
	dbPath, err := resolveDBPath()
	if err != nil {
		fmt.Printf("error resolviendo ruta DB: %v\n", err)
		os.Exit(1)
	}

	db, err := sqlite.NewDB(dbPath, database.Migrations)
	if err != nil {
		fmt.Printf("error inicializando DB: %v\n", err)
		os.Exit(1)
	}

	app := &App{db: db}
	app.wire()
	return app
}

func (a *App) Startup(ctx context.Context) {}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func resolveDBPath() (string, error) {
	// En desarrollo usa la raíz del proyecto
	if os.Getenv("ENV") == "development" {
		return "inventario_dev.db", nil
	}

	// En producción usa el directorio del usuario
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
