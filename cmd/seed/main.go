package main

import (
	"log"
	"time"

	"inventario/app"
	"inventario/database"
	"inventario/internal/platform/sqlite"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbPath, err := app.ResolveDBPath()
	if err != nil {
		log.Fatalf("resolving db path: %v", err)
	}

	db, err := sqlite.NewDB(dbPath, database.Migrations)
	if err != nil {
		log.Fatalf("opening db: %v", err)
	}
	defer db.Close()

	now := time.Now()

	// --- categories ---
	type catRow struct {
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	categories := []catRow{
		{Name: "Electrónica", CreatedAt: now, UpdatedAt: now},
		{Name: "Ropa", CreatedAt: now, UpdatedAt: now},
		{Name: "Alimentos", CreatedAt: now, UpdatedAt: now},
		{Name: "Hogar", CreatedAt: now, UpdatedAt: now},
	}

	catIDs := make(map[string]int64, len(categories))
	for _, c := range categories {
		res, err := db.NamedExec(
			`INSERT OR IGNORE INTO categories (name, created_at, updated_at) VALUES (:name, :created_at, :updated_at)`,
			c,
		)
		if err != nil {
			log.Fatalf("inserting category %q: %v", c.Name, err)
		}
		id, _ := res.LastInsertId()
		if id == 0 {
			// already existed — fetch it
			if err := db.Get(&id, `SELECT id FROM categories WHERE name = ?`, c.Name); err != nil {
				log.Fatalf("fetching category %q: %v", c.Name, err)
			}
		}
		catIDs[c.Name] = id
	}
	log.Printf("categories ready: %v", catIDs)

	// --- products ---
	elec := catIDs["Electrónica"]
	ropa := catIDs["Ropa"]
	alim := catIDs["Alimentos"]
	hogar := catIDs["Hogar"]

	type productRow struct {
		CategoryID *int64    `db:"category_id"`
		Name       string    `db:"name"`
		Code       string    `db:"code"`
		Barcode    string    `db:"barcode"`
		Price      int64     `db:"price"`
		Stock      int       `db:"stock"`
		MinStock   int       `db:"min_stock"`
		Active     bool      `db:"active"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}

	products := []productRow{
		{CategoryID: &elec, Name: "Audífonos Bluetooth", Code: "ELEC-001", Barcode: "7501000000001", Price: 599, Stock: 25, MinStock: 5, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &elec, Name: "Cable USB-C", Code: "ELEC-002", Barcode: "7501000000002", Price: 89, Stock: 100, MinStock: 20, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &elec, Name: "Cargador Inalámbrico", Code: "ELEC-003", Barcode: "7501000000003", Price: 349, Stock: 3, MinStock: 5, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &elec, Name: "Mouse Inalámbrico", Code: "ELEC-004", Barcode: "7501000000004", Price: 279, Stock: 18, MinStock: 5, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &ropa, Name: "Playera Algodón M", Code: "ROPA-001", Barcode: "7502000000001", Price: 150, Stock: 60, MinStock: 10, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &ropa, Name: "Pantalón Jean 32", Code: "ROPA-002", Barcode: "7502000000002", Price: 450, Stock: 30, MinStock: 8, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &ropa, Name: "Sudadera con Capucha L", Code: "ROPA-003", Barcode: "7502000000003", Price: 380, Stock: 2, MinStock: 5, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &alim, Name: "Café Molido 500g", Code: "ALIM-001", Barcode: "7503000000001", Price: 120, Stock: 45, MinStock: 10, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &alim, Name: "Aceite de Oliva 1L", Code: "ALIM-002", Barcode: "7503000000002", Price: 95, Stock: 8, MinStock: 10, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &alim, Name: "Avena Integral 1kg", Code: "ALIM-003", Barcode: "7503000000003", Price: 55, Stock: 70, MinStock: 15, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &hogar, Name: "Lámpara LED de Escritorio", Code: "HOGAR-001", Barcode: "7504000000001", Price: 220, Stock: 12, MinStock: 3, Active: true, CreatedAt: now, UpdatedAt: now},
		{CategoryID: &hogar, Name: "Juego de Sábanas Queen", Code: "HOGAR-002", Barcode: "7504000000002", Price: 650, Stock: 0, MinStock: 3, Active: false, CreatedAt: now, UpdatedAt: now},
		{CategoryID: nil, Name: "Artículo Sin Categoría", Code: "GEN-001", Barcode: "7500000000001", Price: 10, Stock: 50, MinStock: 5, Active: true, CreatedAt: now, UpdatedAt: now},
	}

	inserted := 0
	for _, p := range products {
		_, err := db.NamedExec(
			`INSERT OR IGNORE INTO products
				(category_id, name, code, barcode, price, stock, min_stock, active, created_at, updated_at)
			 VALUES
				(:category_id, :name, :code, :barcode, :price, :stock, :min_stock, :active, :created_at, :updated_at)`,
			p,
		)
		if err != nil {
			log.Printf("skipping %q: %v", p.Name, err)
			continue
		}
		inserted++
	}

	log.Printf("seed done: %d products inserted into %s", inserted, dbPath)
}
