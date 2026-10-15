# InventarioWails

An inventory management desktop application built with Go, Wails, Vue, TypeScript and SQLite. It allows controlling categories, products, stock adjustments and sales in a simple local application.

## Overview

This project follows a layered architecture based on Go and Wails bindings for the desktop shell, with a Vue frontend and a SQLite database. It is designed to manage a small to medium retail inventory flow without requiring a separate server.

## Functional requirements

The application must support the following business capabilities:

- Category management
  - Create, list, update and delete categories.
  - Validate that the category name is required and unique.

- Product management
  - Create, list, search by identifier, update and delete products.
  - Assign a product to a category (optional).
  - Validate required fields such as name, code, price and stock.
  - Prevent duplicate product codes.
  - Toggle product active/inactive state.
  - Detect low stock based on the minimum stock threshold.

- Inventory control
  - Register stock adjustments.
  - Record movement type and optional notes.
  - Validate movement quantities and ensure there is enough stock for reductions.
  - List stock movement history per product.

- Sales management
  - Create sales with multiple products and quantities.
  - Validate that the sale contains at least one item.
  - Prevent invalid or out-of-stock product quantities.
  - Calculate total and line subtotals.
  - List and retrieve sales records.

- Data persistence
  - Store data locally in SQLite.
  - Use migrations to maintain the schema.
  - Keep the app usable in desktop mode without external services.

## Models

### Category

- `ID: int64`
- `Name: string`
- `CreatedAt: time.Time`
- `UpdatedAt: time.Time`

### Product

- `ID: int64`
- `CategoryID: *int64`
- `Name: string`
- `Code: string`
- `Barcode: string`
- `Price: float64`
- `Stock: int`
- `MinStock: int`
- `Active: bool`
- `CreatedAt: time.Time`
- `UpdatedAt: time.Time`

### Sale

- `ID: int64`
- `Total: float64`
- `Items: []SaleItem`
- `CreatedAt: time.Time`

### SaleItem

- `ID: int64`
- `SaleID: int64`
- `ProductID: int64`
- `Quantity: int`
- `UnitPrice: float64`
- `Subtotal: float64`
- `CreatedAt: time.Time`

### StockMovement

- `ID: int64`
- `ProductID: int64`
- `Quantity: int`
- `MovementType: MovementType`
- `Notes: string`
- `CreatedAt: time.Time`

Movement types supported by the domain model:

- `adjustment`
- `sale`
- `return`

## Architecture

The project is organized by responsibilities:

- `app/`: Wails application bootstrap and lifecycle setup.
- `database/`: SQLite migrations and migration loader.
- `internal/domain/`: business entities and validation rules.
- `internal/application/`: use cases for categories, products, stock and sales.
- `internal/handler/`: Wails-exposed handlers used by the frontend.
- `internal/infrastructure/sqlite/`: database access implementation.
- `frontend/`: Vue + TypeScript frontend built with Vite.

## Technology stack

- Go 1.26
- [Wails](https://wails.io/docs/gettingstarted/installation) v2 for desktop application shell
- Vue 3 + TypeScript for the frontend
- Vite for frontend build and dev server
- SQLite for local persistence
- `sqlx` for database access
- `goose` for schema migrations

## Useful commands

Run the application in development mode:

```bash
wails dev
```

Build the desktop application for production:

```bash
wails build
```

Run backend tests:

```bash
go test ./...
```

## Database configuration

The app uses SQLite and resolves the database path automatically:

- Development mode: `inventario_dev.db` in the project root
- Production mode: the user configuration directory under `inventario/inventario.db`

This behavior is handled in `app/app.go`.

## Notes

This project is a local desktop inventory system intended for practical use in small operations. The repository is structured to keep business logic centralized in the domain and application layers, while the desktop shell and database access remain isolated.
