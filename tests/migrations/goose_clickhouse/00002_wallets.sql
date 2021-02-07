-- +goose Up
CREATE TABLE wallets (
    id      BIGINT,
    user_id BIGINT
) engine=MergeTree() order by id;

-- +goose Down
DROP TABLE IF EXISTS wallets;