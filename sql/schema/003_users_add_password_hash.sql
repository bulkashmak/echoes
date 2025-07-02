-- +goose Up
ALTER TABLE users
ADD COLUMN password_hash TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE users
DROP COLUMN password_hash;

