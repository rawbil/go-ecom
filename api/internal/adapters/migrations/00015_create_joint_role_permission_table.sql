-- +goose Up
CREATE TABLE IF NOT EXISTS role_permissions(
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,

    PRIMARY KEY(role_id, permission_id)
);

-- +goose Down
DROP TABLE IF EXISTS role_permissions;
