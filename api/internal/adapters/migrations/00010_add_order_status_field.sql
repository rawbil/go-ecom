-- +goose Up
ALTER TABLE orders
ADD COLUMN order_status ENUM("pending", "checked", "cancelled") DEFAULT "pending";

-- +goose Down
SELECT 'down SQL query';
