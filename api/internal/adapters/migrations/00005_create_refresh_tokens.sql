-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
	refresh_token TEXT NOT NULL,
	user_id BIGINT NOT NULL,
	issued_at DATETIME DEFAULT NOW() NOT NULL,
	expires_at DATETIME NOT NULL,
	CONSTRAINT fk_uid FOREIGN KEY (user_id) REFERENCES users(user_id)
	
);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
