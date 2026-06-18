-- +goose Up
ALTER TABLE products
ADD COLUMN quantity INT NOT NULL;

-- +goose Down
SELECT 'down SQL query';
