CREATE TABLE test_table (
  id BIGSERIAL PRIMARY KEY,
  inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX test_table_idx ON test_table (inserted_at);
