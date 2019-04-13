DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'things') THEN
        create type things AS ENUM ('hello', 'world');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS doge (
  id int,
  th things
);
