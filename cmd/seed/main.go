package main

import (
	"log"
	"time"

	"inventario/app"
	"inventario/database"
	"inventario/internal/platform/gormdb"

	"github.com/joho/godotenv"
)

// Local persistence models: the seed writes straight to the tables and keeps
// no dependency on the modules' unexported repository models.

type seedCategory struct {
	ID        int64 `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (seedCategory) TableName() string { return "categories" }

type seedProduct struct {
	ID         int64 `gorm:"primaryKey"`
	CategoryID *int64
	Name       string
	Code       string
	Barcode    string
	Price      int64
	Stock      int
	MinStock   int
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (seedProduct) TableName() string { return "products" }

func main() {
	_ = godotenv.Load()

	dbPath, err := app.ResolveDBPath()
	if err != nil {
		log.Fatalf("resolving db path: %v", err)
	}

	db, err := gormdb.Open(dbPath)
	if err != nil {
		log.Fatalf("opening db: %v", err)
	}
	defer gormdb.Close(db)

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrating db: %v", err)
	}

	now := time.Now()

	// --- categories ---
	categoryNames := []string{"Electrónica", "Ropa", "Alimentos", "Hogar"}

	catIDs := make(map[string]int64, len(categoryNames))
	for _, name := range categoryNames {
		// Insert or fetch the existing row, keyed on the unique name.
		c := seedCategory{Name: name}
		if err := db.Where(seedCategory{Name: name}).
			Attrs(seedCategory{CreatedAt: now, UpdatedAt: now}).
			FirstOrCreate(&c).Error; err != nil {
			log.Fatalf("seeding category %q: %v", name, err)
		}
		catIDs[name] = c.ID
	}
	log.Printf("categories ready: %v", catIDs)

	// --- products ---
	elec := catIDs["Electrónica"]
	ropa := catIDs["Ropa"]
	alim := catIDs["Alimentos"]
	hogar := catIDs["Hogar"]

	products := []seedProduct{
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
		// Insert only if no product has this code yet (codes are unique).
		row := p
		res := db.Where(seedProduct{Code: p.Code}).Attrs(p).FirstOrCreate(&row)
		if res.Error != nil {
			log.Printf("skipping %q: %v", p.Name, res.Error)
			continue
		}
		if res.RowsAffected > 0 {
			inserted++
		}
	}

	log.Printf("seed done: %d products inserted into %s", inserted, dbPath)
}
