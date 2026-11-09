# InventarioWails

**A lightweight desktop inventory management system built with Go, Wails 3, Vue 3, TypeScript, and SQLite.**

InventarioWails is a local-first desktop application designed to simplify inventory management for small and medium-sized retail operations. It provides tools to manage products, organize categories, track stock movements, and process sales — all without requiring an external server.

## Overview

The application is organized as a set of self-contained business **modules**, each following a clean, layered design that keeps business rules independent of the desktop framework and the database.

- **Local-first:** All data is stored locally using SQLite (pure-Go driver, no CGO).
- **Modular architecture:** Each feature (category, product, stock, sale) is a vertical slice with its own domain, application, handler, and infrastructure layers.
- **Inventory tracking:** Monitor stock levels, adjustments, and product movements.
- **Sales management:** Register sales containing multiple products and quantities.
- **Desktop experience:** A Vue 3 interface integrated with Go through Wails 3 bindings.

## Features

### Categories

- Create, list, search, update, and delete (soft-delete) categories.
- Restore previously deleted categories.
- Validate required category names and prevent duplicates.
- Paginated listing with optional "only inactive" filtering.

### Products

- Create, list, search, update, and delete products.
- Assign products to optional categories.
- Validate required fields, including name, code, price, and stock.
- Prevent duplicate product codes.
- Activate or deactivate products.
- Identify low-stock products based on configurable minimum stock levels.
- Paginated listing.

### Inventory Control

- Register stock adjustments with a movement type and optional notes.
- Validate quantities and prevent stock reductions below available inventory.
- Review stock movement history for an individual product.

### Sales

- Create sales containing multiple products.
- Validate that each sale includes at least one item.
- Prevent invalid quantities and sales exceeding available stock.
- Calculate item subtotals and sale totals.
- List and retrieve sales records.

### Data Persistence

- Store application data locally in SQLite.
- Manage database schema changes through embedded Goose migrations applied automatically at startup.
- Operate independently of external services.

## Domain Models

The application uses the following core entities (Go domain types).

### Category

| Field | Type | Notes |
| --- | --- | --- |
| `ID` | `int64` | |
| `Name` | `string` | Unique |
| `CreatedAt` | `time.Time` | |
| `UpdatedAt` | `time.Time` | |
| `DeletedAt` | `*time.Time` | `nil` when active (soft delete) |

`Category.IsActive()` reports whether `DeletedAt` is `nil`.

### Product

| Field | Type | Notes |
| --- | --- | --- |
| `ID` | `int64` | |
| `CategoryID` | `*int64` | Optional; `NULL` set on category delete |
| `Name` | `string` | |
| `Code` | `string` | Unique |
| `Barcode` | `string` | |
| `Price` | `int64` | Monetary amount stored as an integer |
| `Stock` | `int` | |
| `MinStock` | `int` | |
| `Active` | `bool` | |
| `CreatedAt` | `time.Time` | |
| `UpdatedAt` | `time.Time` | |

`Product.HasLowStock()` returns `true` when `Stock <= MinStock`.

### Sale

| Field | Type |
| --- | --- |
| `ID` | `int64` |
| `Total` | `float64` |
| `Items` | `[]SaleItem` |
| `CreatedAt` | `time.Time` |

### SaleItem

| Field | Type |
| --- | --- |
| `ID` | `int64` |
| `SaleID` | `int64` |
| `ProductID` | `int64` |
| `Quantity` | `int` |
| `UnitPrice` | `float64` |
| `Subtotal` | `float64` |
| `CreatedAt` | `time.Time` |

### Movement (stock)

| Field | Type |
| --- | --- |
| `ID` | `int64` |
| `ProductID` | `int64` |
| `Quantity` | `int` |
| `MovementType` | `MovementType` |
| `Notes` | `string` |
| `CreatedAt` | `time.Time` |

Supported movement types:

- `adjustment` — Manual inventory adjustments.
- `sale` — Stock changes associated with sales.
- `return` — Inventory changes associated with returned products.

## Architecture

The codebase uses a **modular (vertical-slice) architecture**. Each business capability lives under `internal/modules/<feature>` and owns its full stack, while cross-cutting concerns live under `internal/platform`.

```text
InventarioWails/
├── main.go                       # Wails 3 entrypoint; embeds frontend, registers services
├── app/
│   ├── app.go                    # App bootstrap: DB path resolution + lifecycle
│   └── wire.go                   # Dependency wiring across modules
├── cmd/
│   └── seed/                     # Standalone CLI to seed the database
├── database/
│   ├── migration.go              # Embeds migrations/*.sql
│   └── migrations/               # Goose SQL migrations
├── internal/
│   ├── modules/                  # One self-contained module per feature
│   │   ├── category/
│   │   │   ├── category.go       # Composition facade (builds the object graph)
│   │   │   ├── domain/           # Entities, errors, queries, repository interface, validators
│   │   │   ├── application/      # Use cases (create, list, find, update, delete, restore …)
│   │   │   ├── handler/          # Wails-exposed handler + request/response DTOs
│   │   │   └── infrastructure/   # SQLite repository implementation
│   │   ├── product/              # Same internal layout
│   │   ├── stock/                # Same internal layout
│   │   └── sale/                 # Same internal layout
│   └── platform/                 # Shared, domain-agnostic infrastructure
│       ├── sqlite/               # DB connection, migrations runner, query builder, errors
│       ├── pagination/           # Generic pagination params/results
│       └── responses/            # Shared API response envelopes
├── frontend/                     # Vue 3 + TypeScript + Vite app
│   ├── src/                      # Views, components, composables, services, router
│   └── bindings/                 # Auto-generated Wails TypeScript bindings
└── build/                        # Wails 3 build config and per-platform Taskfiles
```

### Layers within each module

| Layer | Responsibility |
| --- | --- |
| `domain` | Entities, invariants, errors, query objects, repository interfaces, and validators. |
| `application` | Use cases that orchestrate domain logic (one type per operation). |
| `handler` | Exposes operations to the frontend as Wails services; maps DTOs ⇆ domain. |
| `infrastructure` | Implements the domain repository interfaces against SQLite. |
| `<feature>.go` | A `New(db)` facade that assembles the module and returns its handler. |

Modules communicate through repository interfaces. For example, [app/wire.go](app/wire.go) injects the product repository into the stock and sale modules so inventory stays consistent across features.

### Wails services and exposed methods

`main.go` registers four handlers as Wails services, generating strongly typed TS bindings:

| Service | Methods |
| --- | --- |
| `CategoryHandler` | `ListCategories`, `ListCategoryOptions`, `GetCategory`, `CreateCategory`, `UpdateCategory`, `DeleteCategory`, `RestoreCategory` |
| `ProductHandler` | `ListProducts`, `GetProduct`, `CreateProduct`, `UpdateProduct`, `DeleteProduct`, `ToggleActiveProduct` |
| `StockHandler` | `AdjustStock`, `ListStockMovements` |
| `SaleHandler` | `CreateSale`, `ListSales`, `GetSale` |

## Technology Stack

### Backend

| Technology | Purpose |
| --- | --- |
| [Go](https://go.dev/) 1.26 | Backend and business logic |
| [Wails 3](https://v3alpha.wails.io/) (`v3.1.0`) | Native desktop application framework |
| [SQLite](https://sqlite.org/) via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Pure-Go local database driver (no CGO) |
| [sqlx](https://github.com/jmoiron/sqlx) | Database access helpers |
| [goose](https://github.com/pressly/goose) | Embedded schema migrations |
| [godotenv](https://github.com/joho/godotenv) | `.env` loading |

### Frontend

| Technology | Purpose |
| --- | --- |
| [Vue 3](https://vuejs.org/) | Frontend framework |
| [TypeScript](https://www.typescriptlang.org/) | Type safety |
| [Vite](https://vite.dev/) | Build tooling and dev server |
| [Vue Router](https://router.vuejs.org/) | Client-side routing |
| [Pinia](https://pinia.vuejs.org/) | State management |
| [PrimeVue](https://primevue.org/) + [PrimeIcons](https://primevue.org/icons/) | UI components |
| [Tailwind CSS](https://tailwindcss.com/) 4 | Styling |
| [pnpm](https://pnpm.io/) | Package manager |

## Getting Started

### Prerequisites

Install the following tools before running the application:

- Go 1.26
- Node.js and [pnpm](https://pnpm.io/)
- [Wails 3 CLI](https://v3alpha.wails.io/getting-started/installation/) (`wails3`)
- Platform build dependencies required by Wails 3

> SQLite uses the pure-Go `modernc.org/sqlite` driver, so a C compiler is **not** required for the database.

### Installation

```bash
git clone <repository-url>
cd inventariowails

cd frontend
pnpm install
cd ..
```

### Development

Run the application in development mode (live reload for Go and frontend):

```bash
wails3 dev
```

### Production Build

Build the desktop application for the host platform:

```bash
wails3 build
```

Package a distributable production build:

```bash
wails3 package
```

Run `wails3 --help` to see every available command.

### Seeding the database

A standalone CLI seeds categories and products into the configured database:

```bash
go run ./cmd/seed
```

### Testing

Run the Go test suite:

```bash
go test ./...
```

## Configuration

### Environment variables

Copy `.env.example` to `.env` and adjust it for your environment (`.env` is gitignored). Both the app and the seed command read these values.

| Variable | Description |
| --- | --- |
| `ENV` | Set to `development` to use `inventario_dev.db` in the working directory. Any other value (or unset) uses the production path. |
| `DB_PATH` | Optional explicit database file path. When set, it takes precedence over `ENV`. |

Database path resolution (see [app/app.go](app/app.go), `ResolveDBPath`):

| Condition | Database location |
| --- | --- |
| `DB_PATH` is set | The value of `DB_PATH` |
| `ENV=development` | `inventario_dev.db` in the working directory |
| Otherwise (production) | `inventario/inventario.db` under the OS user config directory |

The SQLite connection enables foreign keys and WAL journaling, and runs with a single open connection.

### Build configuration

Application metadata (product name, identifier, version, icons) for packaging lives in [build/config.yml](build/config.yml). Per-platform build and packaging steps are defined in the Taskfiles under `build/<platform>/`.

## Design Principles

- **Modularity:** Each feature is a self-contained slice that can be developed and tested in isolation.
- **Separation of concerns:** Within a module, domain, application, handler, and infrastructure each have a single responsibility.
- **Dependency inversion:** Modules depend on repository interfaces, not concrete implementations; wiring happens at the application root.
- **Local persistence:** The application works without a separate backend server.
- **Data integrity:** Domain validation and database constraints prevent invalid inventory and sales records.

## Project Status

InventarioWails is a practical foundation for desktop inventory management. Its modular structure supports incremental development of category management, product catalogues, stock control, and sales workflows.

---

*Built with Go, Wails 3, Vue 3, TypeScript, and SQLite.*
