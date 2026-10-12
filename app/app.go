package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	dbPath, err := resolveDBPath()
	if err != nil {
		fmt.Printf("error resolviendo ruta DB: %v\n", err)
		os.Exit(1)
	}

	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		fmt.Printf("error inicializando DB: %v\n", err)
		os.Exit(1)
	}

	a.db = db

	a.wire()
}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func resolveDBPath() (string, error) {
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
