-- +goose Up
ALTER TABLE users
ADD COLUMN refresh_token_id BIGINT,
ADD CONSTRAINT fk_rt 
FOREIGN KEY (refresh_token_id) REFERENCES refresh_tokens(id);

-- +goose Down
SELECT 'down SQL query';
