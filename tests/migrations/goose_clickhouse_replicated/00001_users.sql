-- +goose Up
CREATE TABLE users_replicated ON CLUSTER '{cluster}' (
    id BIGINT,
    name VARCHAR(128),
    migastas BIGINT DEFAULT 0
) engine=MergeTree() order by id;

-- +goose Down
DROP TABLE IF EXISTS users_replicated ON CLUSTER '{cluster}';
