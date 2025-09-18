CREATE TABLE IF NOT EXISTS test_table (
  id serial PRIMARY KEY,
  msg text,
  created_at timestamptz DEFAULT now()
);

INSERT INTO test_table (msg) VALUES ('hello from master init');
