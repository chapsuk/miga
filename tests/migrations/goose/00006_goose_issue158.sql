-- +goose Up

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'things') THEN
        create type things AS ENUM ('hello', 'world');
    END IF;
END
$$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS doge (
  id int,
  th things
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS doge;
