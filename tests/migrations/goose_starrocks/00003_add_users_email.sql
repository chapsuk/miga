-- +goose Up
ALTER TABLE users ADD COLUMN email STRING;

-- +goose Down
ALTER TABLE users DROP COLUMN email;
