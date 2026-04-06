-- Rollback report indexes
DROP INDEX IF EXISTS idx_transactions_deleted_transacted_type;
DROP INDEX IF EXISTS idx_transactions_deleted_transacted_account;
DROP INDEX IF EXISTS idx_transactions_deleted_transacted_category;
DROP INDEX IF EXISTS idx_transactions_transacted_date;
