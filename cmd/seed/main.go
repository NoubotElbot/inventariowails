package main

import (
	"log"

	"inventario/app"
	"inventario/database"
	"inventario/internal/infrastructure/sqlite"

	"github.com/joho/godotenv"
)

var categories = []string{
	"Confites",
	"Lacteos",
	"Bebidas",
	"Carnes",
	"Verduras",
	"Frutas",
	"Panaderia",
	"Congelados",
	"Abarrotes",
	"Snacks",
	"Productos de limpieza",
	"Higiene personal",
	"Mascotas",
}

func main() {
	_ = godotenv.Load()

	dsn, err := app.ResolveDBPath()
	if err != nil {
		log.Fatalf("resolver ruta de DB: %v", err)
	}

	db, err := sqlite.NewDB(dsn, database.Migrations)
	if err != nil {
		log.Fatalf("abrir DB: %v", err)
	}
	defer db.Close()

	for _, name := range categories {
		_, err := db.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, name)
		if err != nil {
			log.Fatalf("insertar categoría %q: %v", name, err)
		}
	}

	log.Printf("seed completado: %d categorías insertadas", len(categories))
}
