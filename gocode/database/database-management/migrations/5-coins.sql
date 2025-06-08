CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('crypto', 'fiat', 'stablecoin')),
    pegged_to INTEGER REFERENCES currencies(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO currencies (name, type) VALUES 
  ('BTC', 'crypto'),
  ('ETH', 'crypto'),
  ('FTM', 'crypto'),
  ('SOL', 'crypto'),
  ('ADA', 'crypto'),
  ('DOT', 'crypto'),
  ('SHIB', 'crypto'),
  ('DOGE', 'crypto'),
  ('ROSE', 'crypto'),
  ('POL', 'crypto'),
  ('MATIC', 'crypto'),
  ('NEAR', 'crypto'),
  ('AVAX', 'crypto');


INSERT INTO currencies (name, type) VALUES 
('USD', 'fiat'),
('GBP', 'fiat'),
('EUR', 'fiat');


INSERT INTO currencies (name, type, pegged_to) VALUES 
  ('USDT', 'stablecoin', (SELECT id FROM currencies WHERE name = 'USD')),
  ('BUSD', 'stablecoin', (SELECT id FROM currencies WHERE name = 'USD')),
  ('USDC', 'stablecoin', (SELECT id FROM currencies WHERE name = 'USD'));