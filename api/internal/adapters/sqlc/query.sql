
-- name: ListUsers :many
SELECT * FROM users
WHERE (
    (sqlc.arg(username) = '' OR username LIKE CONCAT('%', sqlc.arg(username),'%')) 
    AND (sqlc.arg(email) = '' OR email LIKE CONCAT('%', sqlc.arg(email), '%')) 
    AND (sqlc.arg(role) = '' OR role LIKE CONCAT('%', sqlc.arg(role), '%'))
)
ORDER BY updated_at DESC
LIMIT ?
OFFSET ?;

-- name: ListUser :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUserById :one
SELECT * FROM users WHERE user_id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password)
VALUES (?, ?, ?);

-- name: UpdateUserEmail :execresult
UPDATE users
SET email=?
WHERE user_id=?;

-- name: UpdateUsername :execresult
UPDATE users
SET username=?
WHERE user_id=?;

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

-- name: UpdateProduct :execresult
UPDATE products
SET price = ?,
    quantity = ?
WHERE product_id = ?;

-- name: ListProduct :one
SELECT * FROM products WHERE product_id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE (
    (sqlc.arg(name) = '' OR product_name LIKE CONCAT('%', sqlc.arg(name), '%')) AND
    (sqlc.arg(min_price) = 0 OR price >= sqlc.arg(min_price)) AND
    (sqlc.arg(max_price) = 0 OR price <= sqlc.arg(max_price))
)
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

-- name: DeleteOrderItemByProductId :exec
DELETE FROM order_items
WHERE product_id = ?;

-- name: AllUsersOrderDetails :many
SELECT orders.order_id, orders.order_status as order_status, orders.created_at, order_items.total_price, order_items.quantity AS order_quantity, product_name, products.price AS product_price, products.quantity AS available_products, username, email FROM orders
INNER JOIN order_items
INNER JOIN products
INNER JOIN users
WHERE orders.order_id = order_items.order_id AND order_items.product_id = products.product_id AND orders.user_id = users.user_id
ORDER BY orders.created_at DESC;

-- name: IdUserOrderDetails :many
SELECT orders.order_id, orders.order_status as order_status, orders.created_at, order_items.total_price, order_items.quantity AS order_quantity, product_name, products.price AS product_price, products.quantity AS available_products, username, email FROM orders
INNER JOIN order_items
INNER JOIN products
INNER JOIN users
WHERE orders.order_id = order_items.order_id AND order_items.product_id = products.product_id AND orders.user_id = users.user_id AND users.user_id = ?
ORDER BY orders.created_at;

-- name: GetOrderById :one
SELECT orders.order_id, orders.order_status as order_status, orders.created_at, order_items.total_price, order_items.quantity AS order_quantity, product_name, products.price AS product_price, products.quantity AS available_products, username, email FROM orders
INNER JOIN order_items
INNER JOIN products
INNER JOIN users
WHERE orders.order_id = order_items.order_id AND order_items.product_id = products.product_id AND orders.user_id = users.user_id AND users.user_id = ? AND orders.order_id = ?;

-- name: GetPaidProductOrders :many
SELECT orders.order_id, orders.user_id, orders.created_at, orders.order_status, order_items.product_id 
FROM orders
INNER JOIN order_items
ON orders.order_id = order_items.order_id 
WHERE order_items.product_id = ? AND orders.order_status = "paid";

-- name: CancelOrder :execresult
UPDATE orders
SET order_status = "cancelled"
WHERE order_id = ? AND order_status = "pending";

-- name: DeleteOrderItemByProductID :exec
DELETE FROM order_items
WHERE product_id = ?;

-- name: DeleteOrderswithoutItems :exec
DELETE FROM orders
WHERE NOT EXISTS(
    SELECT 1
    FROM order_items
    WHERE order_items.order_id = orders.order_id
);

-- name: GetPermissionID :one
SELECT id FROM user_permissions
WHERE permission = ?;

-- name: GetRoleID :one
SELECT id FROM user_roles
WHERE role = ?;

-- name: CreateUserRole :execresult
INSERT INTO roles(user_id, role_id)
VALUES (?, ?);

-- name: GetUserPermissions :many
SELECT up.permission FROM user_permissions up
JOIN role_permissions rp
ON up.id = rp.permission_id
JOIN roles ur
ON ur.role_id = rp.role_id
WHERE ur.user_id = ?;