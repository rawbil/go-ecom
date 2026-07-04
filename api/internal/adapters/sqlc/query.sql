
-- name: ListUsers :many
SELECT * FROM users
WHERE (
    (sqlc.arg(username) = '' OR username LIKE CONCAT('%', sqlc.arg(username),'%')) 
    AND (sqlc.arg(email) = '' OR email LIKE CONCAT('%', sqlc.arg(email), '%')) 
    AND (sqlc.arg(role) OR role LIKE CONCAT('%', sqlc.arg(role), '%'))
)
ORDER BY updated_at
LIMIT ?
OFFSET ?;

-- name: ListUser :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUserById :one
SELECT * FROM users WHERE user_id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password)
VALUES (?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users WHERE email = ?;


-- name: UpdatePassword :execresult
UPDATE users
SET password = ?
WHERE user_id = ?;

-- name: CreateRefreshToken :execresult
INSERT INTO refresh_tokens (refresh_token, user_id, issued_at, expires_at)
VALUES (?, ?, ?, ?);

-- name: UpdateRefreshToken :execresult
UPDATE refresh_tokens
SET refresh_token = ?,
    issued_at = ?,
    expires_at = ?
WHERE user_id = ?;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id = ?;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE user_id = ?;


-- name: CreateProduct :execresult
INSERT INTO products(product_name, price, quantity)
VALUES (?, ?, ?);

-- name: ListProduct :one
SELECT * FROM products WHERE product_id = ? LIMIT 1;

-- name: ListProducts :many
SELECT *
FROM products
WHERE
    (sqlc.arg(name) = '' OR product_name LIKE CONCAT('%', sqlc.arg(name), '%'))
    AND (sqlc.arg(min_price) = 0 OR price >= sqlc.arg(min_price))
    AND (sqlc.arg(max_price) = 0 OR price <= sqlc.arg(max_price))
ORDER BY updated_at DESC
LIMIT ?
OFFSET ?;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = ?;

-- name: UpdateProductQuantity :execresult
UPDATE products
SET quantity = ?
WHERE product_id = ?;

-- name: CreateOrder :execresult
INSERT INTO orders(user_id)
VALUES (?);

-- name: CreateOrderItem :execresult
INSERT INTO order_items (order_id, product_id, quantity, total_price)
VALUES (?, ?, ?, ?);

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: ListOrderItems :many
SELECT * FROM order_items
ORDER BY created_at DESC;

-- name: ListOrder :one
SELECT * FROM orders
WHERE order_id = ?;

-- name: ListOrderItem :one
SELECT * FROM order_items
WHERE id = ?;