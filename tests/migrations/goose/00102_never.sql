-- +goose Up
CREATE TABLE never (
    id INT PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS never;