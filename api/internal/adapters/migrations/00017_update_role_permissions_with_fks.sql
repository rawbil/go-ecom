-- +goose Up
ALTER TABLE role_permissions
ADD CONSTRAINT fk_role FOREIGN KEY(role_id) REFERENCES user_roles(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_permission FOREIGN KEY (permission_id) REFERENCES user_permissions(id) ON DELETE CASCADE;

-- +goose Down
SELECT 'down SQL query';
