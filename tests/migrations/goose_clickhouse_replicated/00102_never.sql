-- +goose Up
CREATE TABLE never_replicated ON CLUSTER '{cluster}' (
    id INT
) engine=Memory;

-- +goose Down
DROP TABLE IF EXISTS never_replicated ON CLUSTER '{cluster}';
