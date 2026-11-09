-- +goose Up
ALTER TABLE categories ADD COLUMN deleted_at DATETIME;

-- +goose Down
ALTER TABLE categories DROP COLUMN deleted_at;
