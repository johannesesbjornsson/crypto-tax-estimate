CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    date DATE NOT NULL,
    description VARCHAR(100),
    source VARCHAR(50) NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('Income', 'Buy', 'Sell', 'Lost')),
    amount NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    asset TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);

CREATE TRIGGER set_transaction_timestamp
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)
INSERT INTO transactions (date, description,  source, type, amount, price, asset, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE - INTERVAL '1 day',
  'Salary for May',
  'Manual',
  'Income',
  5000.00,
  1.00,
  'USD',
  id,
  now(),
  now()
FROM target_user;


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)
INSERT INTO transactions (date, description, source, type, amount, price, asset, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE,
  'Bought BTC',
  'Manual',
  'Buy',
  1000.00,
  69420.00,
  'BTC',
  id,
  now(),
  now()
FROM target_user;