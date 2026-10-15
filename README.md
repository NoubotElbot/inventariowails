# InventarioWails

**A lightweight desktop inventory management system built with Go, Wails, Vue, TypeScript, and SQLite.**

InventarioWails is a local-first desktop application designed to simplify inventory management for small and medium-sized retail operations. It provides tools to manage products, organize categories, track stock movements, and process sales — all without requiring an external server.

## Overview

The application follows a layered architecture that separates business logic, application use cases, database access, and the desktop interface.

- **Local-first:** All data is stored locally using SQLite.
- **Modular architecture:** Business rules and infrastructure remain separated.
- **Inventory tracking:** Monitor stock levels, adjustments, and product movements.
- **Sales management:** Register sales containing multiple products and quantities.
- **Desktop experience:** A Vue-powered interface integrated with Go through Wails.

## Features

### Categories

- Create, list, update, and delete categories.
- Validate required category names.
- Prevent duplicate category names.

### Products

- Create, list, search, update, and delete products.
- Assign products to optional categories.
- Validate required fields, including name, code, price, and stock.
- Prevent duplicate product codes.
- Activate or deactivate products.
- Identify low-stock products based on configurable minimum stock levels.

### Inventory Control

- Register stock adjustments.
- Track movement types and optional notes.
- Validate quantities and prevent stock reductions below available inventory.
- Review stock movement history for individual products.

### Sales

- Create sales containing multiple products.
- Validate that each sale includes at least one item.
- Prevent invalid quantities and sales exceeding available stock.
- Calculate item subtotals and sale totals.
- List and retrieve sales records.

### Data Persistence

- Store application data locally in SQLite.
- Manage database schema changes through migrations.
- Operate independently of external services.

## Domain Models

The application uses the following core entities.

### Category

| Field | Type |
| --- | --- |
| `ID` | `int64` |
| `Name` | `string` |
| `CreatedAt` | `time.Time` |
| `UpdatedAt` | `time.Time` |

### Product

| Field | Type |
| --- | --- |
| `ID` | `int64` |
| `CategoryID` | `*int64` |
| `Name` | `string` |
| `Code` | `string` |
| `Barcode` | `string` |
| `Price` | `float64` |
| `Stock` | `int` |
| `MinStock` | `int` |
| `Active` | `bool` |
| `CreatedAt` | `time.Time` |
| `UpdatedAt` | `time.Time` |

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

### StockMovement

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

The project is organized by responsibility, keeping business rules independent of the desktop framework and persistence implementation.

```text
InventarioWails/
├── app/                          # Wails bootstrap and lifecycle
├── database/                     # SQLite migrations and loader
├── internal/
│   ├── domain/                   # Entities and business rules
│   ├── application/              # Application use cases
│   ├── handler/                  # Wails-exposed handlers
│   └── infrastructure/
│       └── sqlite/               # SQLite data access
└── frontend/                     # Vue, TypeScript, and Vite
```

| Layer | Responsibility |
| --- | --- |
| `domain` | Defines entities, invariants, and business rules. |
| `application` | Coordinates use cases for categories, products, inventory, and sales. |
| `handler` | Exposes application operations to the frontend through Wails bindings. |
| `infrastructure/sqlite` | Implements database access using SQLite. |
| `database` | Manages schema migrations and database initialization. |
| `app` | Configures the application lifecycle and dependencies. |
| `frontend` | Provides the desktop user interface. |

## Technology Stack

| Technology | Purpose |
| --- | --- |
| [Go](https://go.dev/) 1.26 | Backend and business logic |
| [Wails](https://wails.io/) v2 | Native desktop application shell |
| [Vue 3](https://vuejs.org/) | Frontend framework |
| [TypeScript](https://www.typescriptlang.org/) | Frontend type safety |
| [Vite](https://vite.dev/) | Frontend tooling and development server |
| [SQLite](https://sqlite.org/) | Local relational database |
| [sqlx](https://github.com/jmoiron/sqlx) | Database access |
| [goose](https://github.com/pressly/goose) | Database schema migrations |

## Getting Started

### Prerequisites

Install the following tools before running the application:

- Go 1.26
- Node.js and pnpm
- Wails CLI v2
- A C compiler and the platform dependencies required by Wails

Refer to the [Wails installation guide](https://wails.io/docs/gettingstarted/installation) for platform-specific setup instructions.

### Installation

Clone the repository and install the frontend dependencies:

```bash
git clone <repository-url>
cd InventarioWails

cd frontend
pnpm install
cd ..
```

### Development

Start the application in development mode:

```bash
wails dev
```

### Production Build

Build the desktop application:

```bash
wails build
```

### Testing

Run the Go test suite:

```bash
go test ./internal/domain/... ./internal/application/...
```

## Database Configuration

The application uses SQLite for local persistence. The database path is resolved automatically by `app/app.go`.

| Environment | Database location |
| --- | --- |
| Development | `inventario_dev.db` in the project root |
| Production | User configuration directory under `inventario/inventario.db` |

Database migrations are managed with Goose to keep the schema versioned and consistent across installations.

## Design Principles

- **Separation of concerns:** Each layer has a clearly defined responsibility.
- **Domain-driven business logic:** Validation and inventory rules are centralized in the domain and application layers.
- **Local persistence:** The application works without a separate backend server.
- **Maintainability:** Modular components make individual features easier to test and extend.
- **Data integrity:** Validation and database constraints help prevent invalid inventory and sales records.

## Project Status

InventarioWails is intended as a practical foundation for desktop inventory management. Its modular structure supports incremental development of category management, product catalogues, stock control, and sales workflows.

---

*Built with Go, Wails, Vue, TypeScript, and SQLite.*
