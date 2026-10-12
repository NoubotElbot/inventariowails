-- +goose Up
CREATE TABLE
    IF NOT EXISTS sale_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sale_id INTEGER NOT NULL,
        product_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL CHECK (quantity > 0),
        unit_price INTEGER NOT NULL CHECK (unit_price >= 0),
        subtotal INTEGER NOT NULL CHECK (subtotal >= 0),
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sale_id) REFERENCES sales (id),
        FOREIGN KEY (product_id) REFERENCES products (id)
    );

-- +goose Down
DROP TABLE IF EXISTS sale_items;