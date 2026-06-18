-- +goose Up
CREATE TABLE IF NOT EXISTS products(
    product_id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    product_name TEXT NOT NULL,
    price BIGINT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT NOW(),
    updated_at DATETIME NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_price CHECK(price > 0)
);

-- +goose Down
DROP TABLE IF EXISTS products;
