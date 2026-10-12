-- +goose Up
CREATE TABLE
    IF NOT EXISTS stock_movements (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        product_id INT NOT NULL,
        quantity INT NOT NULL CHECK (quantity >= 0),
        notes TEXT DEFAULT NULL,
        movement_type VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (product_id) REFERENCES products (id)
    );

-- +goose Down
DROP TABLE stock_movements;