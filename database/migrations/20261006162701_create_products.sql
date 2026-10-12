-- +goose Up
CREATE TABLE
    IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER REFERENCES categories (id) ON DELETE SET NULL,
        name TEXT NOT NULL,
        code TEXT NOT NULL UNIQUE,
        barcode TEXT,
        price REAL NOT NULL DEFAULT 0,
        stock INTEGER NOT NULL DEFAULT 0,
        min_stock INTEGER NOT NULL DEFAULT 0,
        active INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS products;