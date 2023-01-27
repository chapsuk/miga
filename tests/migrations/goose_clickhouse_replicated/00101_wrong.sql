-- +goose Up
ALTER TABLE users_replicated ON CLUSTER '{cluster}' ADD COLUMN email VARCHAR;

-- +goose Down
ALTER TABLE users_replicated ON CLUSTER '{cluster}' DROP COLUMN IF EXISTS email;
