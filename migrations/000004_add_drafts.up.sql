CREATE TABLE IF NOT EXISTS drafts (
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
    notes TEXT,
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'rejected')),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_drafts_sub_id ON drafts(sub_id);
CREATE INDEX idx_drafts_account_id ON drafts(account_id);
CREATE INDEX idx_drafts_transfer_account_id ON drafts(transfer_account_id);
CREATE INDEX idx_drafts_category_id ON drafts(category_id);
CREATE INDEX idx_drafts_source ON drafts(source);
CREATE INDEX idx_drafts_status ON drafts(status);
CREATE INDEX idx_drafts_deleted_at ON drafts(deleted_at);