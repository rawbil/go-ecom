-- +goose Up
CREATE TABLE IF NOT EXISTS user_permissions(
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    permission VARCHAR(100) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS user_permissions;
