-- +goose Up
ALTER TABLE users
MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT("user");

-- +goose Down
SELECT 'down SQL query';
