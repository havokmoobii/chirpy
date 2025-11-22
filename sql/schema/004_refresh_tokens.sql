-- +goose Up
CREATE TABLE refresh_tokens (
    token text PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (user_id)
    REFERENCES users(id),
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);

-- +goose Down
DROP TABLE refresh_tokens;