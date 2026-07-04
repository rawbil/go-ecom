-- +goose Up
ALTER TABLE products
DROP CHECK chk_price,
DROP CHECK chk_qty,
ADD CONSTRAINT chk_price CHECK(price >= 0),
ADD CONSTRAINT chk_qty CHECK(quantity >= 0);

-- +goose Down
SELECT 'down SQL query';
