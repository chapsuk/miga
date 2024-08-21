-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGINT,
    name STRING,
    migastas BIGINT DEFAULT "0"
)
PRIMARY KEY (id)
DISTRIBUTED BY HASH (id)
ORDER BY (id);

-- +goose Down
DROP TABLE IF EXISTS users;
