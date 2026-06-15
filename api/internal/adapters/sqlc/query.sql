
-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: ListUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (username, email, password)
VALUES (?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;