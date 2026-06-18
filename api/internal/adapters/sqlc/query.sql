
-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: ListUser :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password)
VALUES (?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users WHERE email = ?;

-- name: CreateProduct :execresult
INSERT INTO products (product_name, price)
VALUES (?, ?);

-- name: ListProduct :one
SELECT * FROM products WHERE product_id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products 
ORDER BY updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = ?;