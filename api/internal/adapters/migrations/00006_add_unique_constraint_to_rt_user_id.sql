-- +goose Up
ALTER TABLE refresh_tokens
MODIFY COLUMN user_id BIGINT NOT NULL UNIQUE;

-- +goose Down
SELECT 'down SQL query';
