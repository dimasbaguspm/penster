-- Composite indexes for report query optimization

-- High-value covering index for GetTotalsByType and other date-range queries
CREATE INDEX idx_transactions_deleted_transacted_type ON transactions(deleted_at, transacted_at, transaction_type);

-- For GetAccountBreakdown (covering index includes account_id)
CREATE INDEX idx_transactions_deleted_transacted_account ON transactions(deleted_at, transacted_at, account_id);

-- For GetCategoryBreakdown (covering index includes category_id)
CREATE INDEX idx_transactions_deleted_transacted_category ON transactions(deleted_at, transacted_at, category_id);

-- For GetTrends - functional index on DATE(transacted_at)
CREATE INDEX idx_transactions_transacted_date ON transactions((DATE(transacted_at)));
