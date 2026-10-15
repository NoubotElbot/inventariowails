package main

import (
	"embed"

	"inventario/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	a := app.NewApp()

	err := wails.Run(&options.App{
		Title:  "Inventario",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  a.Startup,
		OnShutdown: a.Shutdown,
		Bind: []interface{}{
			a.CategoryHandler,
			a.ProductHandler,
			a.StockHandler,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
