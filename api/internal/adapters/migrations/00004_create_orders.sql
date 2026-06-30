-- +goose Up
CREATE TABLE IF NOT EXISTS orders(
    order_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id BIGINT NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    created_at DATETIME NOT NULL DEFAUlT NOW()
);

CREATE TABLE IF NOT EXISTS cart(
    order_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id BIGINT NOT NULL,

    CONSTRAINT fk_cart_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    created_at DATETIME NOT NULL DEFAUlT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    quantity INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(order_id),
    CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(product_id)
);

CREATE TABLE IF NOT EXISTS cart_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    quantity INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_cart_order_id FOREIGN KEY (order_id) REFERENCES orders(order_id),
    CONSTRAINT fk_cart_product_id FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- +goose Down
SELECT 'down SQL query';
