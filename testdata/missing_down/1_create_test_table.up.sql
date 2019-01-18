CREATE TABLE test_table (
  id BIGSERIAL PRIMARY KEY,
  inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON test_table (inserted_at);
