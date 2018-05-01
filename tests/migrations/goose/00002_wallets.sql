-- +goose Up
CREATE TABLE wallets (
    id INT PRIMARY KEY,
    user_id int
);

-- +goose Down
DROP TABLE IF EXISTS wallets;