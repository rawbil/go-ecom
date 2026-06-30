
-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: ListUser :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUserById :one
SELECT * FROM users WHERE user_id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password)
VALUES (?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users WHERE email = ?;

-- name: UpdateUserToken :execresult
UPDATE users
SET refresh_token_id = ?
WHERE user_id = ?;

-- name: CreateRefreshToken :execresult
INSERT INTO refresh_tokens (refresh_token, user_id, issued_at, expires_at)
VALUES (?, ?, ?, ?);

-- name: UpdateRefreshToken :execresult
UPDATE refresh_tokens
SET refresh_token = ?,
    issued_at = ?,
    expires_at = ?
WHERE user_id = ? AND id = ?;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE user_id = ?;

-- name: ListProduct :one
SELECT * FROM products WHERE product_id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products 
ORDER BY updated_at;

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