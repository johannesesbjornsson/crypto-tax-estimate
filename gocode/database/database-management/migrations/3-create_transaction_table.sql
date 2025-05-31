CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    date DATE NOT NULL,
    description VARCHAR(100),
    venue TEXT NOT NULL,
    source VARCHAR(50) NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('Income', 'Buy', 'Sell', 'Lost')),
    amount NUMERIC NOT NULL,
    asset TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER set_transaction_timestamp
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)
INSERT INTO transactions (date, description, venue, source, type, amount, asset, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE - INTERVAL '1 day',
  'Salary for May',
  'Talos',
  'Manual',
  'Income',
  5000.00,
  'USD',
  id,
  now(),
  now()
FROM target_user;


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)
INSERT INTO transactions (date, description, venue, source, type, amount, asset, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE,
  'Bought BTC',
  'Binance',
  'Manual',
  'Buy',
  1000.00,
  'BTC',
  id,
  now(),
  now()
FROM target_user;