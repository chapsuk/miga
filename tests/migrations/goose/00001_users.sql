-- +goose Up
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(128)
);

-- +goose Down
DROP TABLE IF EXISTS users;