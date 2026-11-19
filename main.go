package main

import (
	"embed"
	"io/fs"
	"log"
	"time"

	"inventario/app"
	systemhandler "inventario/internal/modules/system/adapters/handler"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var embeddedAssets embed.FS

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
	application.RegisterEvent[systemhandler.UpdateProgress](systemhandler.UpdateProgressEvent)
}

func main() {
	_ = godotenv.Load()

	assets, err := fs.Sub(embeddedAssets, "frontend/dist")
	if err != nil {
		log.Fatalf("error loading frontend assets: %v", err)
	}

	inventory, err := app.NewApp()
	if err != nil {
		log.Fatalf("error inicializando la aplicación: %v", err)
	}
	defer inventory.Close()

	wailsApp := application.New(application.Options{
		Name: "Inventario",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(inventory.CategoryHandler),
			application.NewService(inventory.ProductHandler),
			application.NewService(inventory.StockHandler),
			application.NewService(inventory.SaleHandler),
			application.NewService(inventory.SystemHandler),
		},
	})

	wailsApp.Window.NewWithOptions(
		application.WebviewWindowOptions{
			Title:  "Inventario",
			Width:  1280,
			Height: 800,
		},
	)

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			wailsApp.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
