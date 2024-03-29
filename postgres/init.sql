CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    hash TEXT,
    status BIGINT,
    bloc_hash TEXT,
    bloc_number BIGINT,
    sender TEXT,
    recipient TEXT,
    address TEXT,
    logs_count SMALLINT,
    input TEXT,
    value TEXT
);