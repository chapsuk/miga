-- +goose Up
SET replication_alter_partitions_sync = 2;
ALTER TABLE users_replicated ON CLUSTER '{cluster}' ADD COLUMN email VARCHAR(128);

-- +goose Down
SET replication_alter_partitions_sync = 2;
ALTER TABLE users_replicated ON CLUSTER '{cluster}' DROP COLUMN email;
