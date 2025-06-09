CREATE TABLE market_prices (
    id SERIAL PRIMARY KEY,
    base_currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    quote_currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    price NUMERIC NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    UNIQUE (base_currency_id, quote_currency_id, timestamp)
);