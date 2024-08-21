-- +goose Up
CREATE TABLE IF NOT EXISTS wallets (
    id      BIGINT,
    user_id BIGINT
)
PRIMARY KEY (id)
DISTRIBUTED BY HASH (id);

-- +goose Down
DROP TABLE IF EXISTS wallets;
