-- +goose Up
ALTER TABLE users
ADD COLUMN role VARCHAR(20) DEFAULT("user");

-- +goose Down
SELECT 'down SQL query';
