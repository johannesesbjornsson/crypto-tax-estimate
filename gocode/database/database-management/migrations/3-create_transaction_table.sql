CREATE TABLE IF NOT EXISTS simple_transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    date DATE NOT NULL,
    description VARCHAR(100),
    source VARCHAR(50) NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('income', 'lost')),
    amount NUMERIC NOT NULL,
    asset TEXT NOT NULL,
    external_id VARCHAR(100),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_stimple_transactions_user_id ON simple_transactions(user_id);

CREATE TRIGGER set_simple_transaction_timestamp
BEFORE UPDATE ON simple_transactions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS trade_transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    date DATE NOT NULL,
    description VARCHAR(100),
    source VARCHAR(50) NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('buy', 'sell')),
    amount NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    asset TEXT NOT NULL,
    quote_currency TEXT NOT NULL,
    external_id VARCHAR(100),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_trade_transactions_user_id ON trade_transactions(user_id);

CREATE TRIGGER set_trade_transactions_timestamp
BEFORE UPDATE ON trade_transactions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)

INSERT INTO simple_transactions (date, description,  source, type, amount, asset, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE - INTERVAL '1 day',
  'Salary for May',
  'Manual',
  'income',
  5000.00,
  'USDT',
  id,
  now(),
  now()
FROM target_user;


WITH target_user AS (
  SELECT id FROM users WHERE email = 'johannes.esbjornsson@gmail.com' LIMIT 1
)
INSERT INTO trade_transactions (date, description, source, type, amount, price, asset, quote_currency, user_id, created_at, updated_at)
SELECT 
  CURRENT_DATE,
  'Bought BTC',
  'Manual',
  'buy',
  1000.00,
  69420.00,
  'BTC',
  'USDT',
  id,
  now(),
  now()
FROM target_user;