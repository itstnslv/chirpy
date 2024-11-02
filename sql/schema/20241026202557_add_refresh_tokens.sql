-- +goose Up
CREATE TABLE refresh_tokens
(
    token      VARCHAR(255) NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);


-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
