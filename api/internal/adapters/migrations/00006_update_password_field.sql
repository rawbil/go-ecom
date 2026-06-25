-- +goose Up
ALTER TABLE users
MODIFY COLUMN password TEXT NOT NULL;

-- +goose Down
SELECT 'down SQL query';
