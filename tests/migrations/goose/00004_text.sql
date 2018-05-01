-- +goose Up
INSERT INTO users (id, name) VALUES (1, 'Abib;Rabib');

-- +goose Down
DELETE FROM users WHERE id = 1;