-- +goose Up
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(128),
    migastas INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS users;