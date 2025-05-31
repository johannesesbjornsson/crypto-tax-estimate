CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    tax_start_day INT CHECK (tax_start_day BETWEEN 1 AND 31),
    tax_start_month INT CHECK (tax_start_month BETWEEN 1 AND 12),
    tax_end_day INT CHECK (tax_end_day BETWEEN 1 AND 31),
    tax_end_month INT CHECK (tax_end_month BETWEEN 1 AND 12),
    currency VARCHAR(4) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

INSERT INTO users (email, name, tax_start_day, tax_start_month, tax_end_day, tax_end_month, currency, created_at, updated_at)
VALUES ('johannes.esbjornsson@gmail.com', 'Johannes', 1, 1, 31, 12, 'GBP', now(), now());