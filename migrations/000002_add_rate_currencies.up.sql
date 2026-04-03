-- Create rate_currencies table
CREATE TABLE IF NOT EXISTS rate_currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_currency TEXT NOT NULL,
    to_currency TEXT NOT NULL,
    rate NUMERIC(18,8) NOT NULL,
    rate_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(from_currency, to_currency, rate_date)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_rate_currencies_currencies ON rate_currencies(from_currency, to_currency);
CREATE INDEX IF NOT EXISTS idx_rate_currencies_rate_date ON rate_currencies(rate_date);
