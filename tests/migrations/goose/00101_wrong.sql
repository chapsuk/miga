-- +goose Up
ALTER TABLE users ADD COLUMN email VARCHAR;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS email;

