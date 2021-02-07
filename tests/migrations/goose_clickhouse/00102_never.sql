-- +goose Up
CREATE TABLE never (
    id INT
) engine=Memory;

-- +goose Down
DROP TABLE IF EXISTS never;
