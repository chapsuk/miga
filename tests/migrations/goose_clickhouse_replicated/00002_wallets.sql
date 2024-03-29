-- +goose Up
CREATE TABLE wallets_replicated ON CLUSTER '{cluster}' (
    id      BIGINT,
    user_id BIGINT
) engine=ReplicatedMergeTree() order by id;

-- +goose Down
DROP TABLE IF EXISTS wallets_replicated ON CLUSTER '{cluster}' SYNC;
