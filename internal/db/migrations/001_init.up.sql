CREATE TABLE IF NOT EXISTS schema_version (
    id         SERIAL PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
