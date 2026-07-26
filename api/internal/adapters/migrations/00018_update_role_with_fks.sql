-- +goose Up
ALTER TABLE roles
ADD CONSTRAINT FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
ADD CONSTRAINT FOREIGN KEY (role_id) REFERENCES user_roles(id) ON DELETE CASCADE;

-- +goose Down
SELECT 'down SQL query';
