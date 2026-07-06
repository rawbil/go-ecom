-- +goose Up
ALTER TABLE orders
MODIFY COLUMN order_status ENUM("pending", "paid", "completed", "cancelled") NOT NULL DEFAULT "pending";

-- +goose Down
SELECT 'down SQL query';
