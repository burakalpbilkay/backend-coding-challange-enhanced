CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMPTZ
);

COPY users(id, name, created_at) FROM '/docker-entrypoint-initdb.d/users.csv' WITH  (FORMAT csv, HEADER true);
