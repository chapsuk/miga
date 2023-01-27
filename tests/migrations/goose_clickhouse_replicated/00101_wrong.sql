-- +goose Up
SET replication_alter_partitions_sync = 2;
ALTER TABLE users_replicated ON CLUSTER '{cluster}' ADD COLUMN email VARCHAR;

-- +goose Down
SET replication_alter_partitions_sync = 2;
ALTER TABLE users_replicated ON CLUSTER '{cluster}' DROP COLUMN IF EXISTS email;
