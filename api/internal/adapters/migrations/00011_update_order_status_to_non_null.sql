-- +goose Up
ALTER TABLE orders
MODIFY COLUMN order_status ENUM("pending", "checked", "cancelled") NOT NULL DEFAULT "pending";

-- +goose Down
SELECT 'down SQL query';
