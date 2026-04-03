-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sub_id UUID NOT NULL DEFAULT gen_random_uuid(),
    account_id INT NOT NULL REFERENCES accounts(id),
    transfer_account_id INT REFERENCES accounts(id),
    category_id INT REFERENCES categories(id),
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('expense', 'income', 'transfer')),
    title TEXT NOT NULL,
    base_amount BIGINT NOT NULL,
    enhanced_amount BIGINT,
    currency TEXT NOT NULL,
    currency_rate NUMERIC(18,8) DEFAULT 1,
    transacted_at DATE NOT NULL,
    notes TEXT,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transactions_sub_id ON transactions(sub_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_transfer_account_id ON transactions(transfer_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_transacted_at ON transactions(transacted_at);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);
