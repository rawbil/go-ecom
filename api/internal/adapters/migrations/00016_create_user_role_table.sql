-- +goose Up
CREATE TABLE IF NOT EXISTS roles( -- allow multiple roles
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,

    PRIMARY KEY (user_id, role_id)
);

-- +goose Down
DROP TABLE IF EXISTS roles ;
