-- name: GetTotalBalance :one
SELECT COALESCE(SUM(balance), 0) as total_balance
FROM accounts
WHERE deleted_at IS NULL;

-- name: GetTotalsByType :many
SELECT
    transaction_type,
    COALESCE(SUM(enhanced_amount), 0) as total
FROM transactions
WHERE deleted_at IS NULL
    AND transacted_at >= @start_date
    AND transacted_at <= @end_date
GROUP BY transaction_type;

-- name: GetCategoryBreakdown :many
SELECT
    c.sub_id as category_id,
    c.name as category_name,
    c.type as category_type,
    COALESCE(SUM(t.enhanced_amount), 0) as total
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id AND c.deleted_at IS NULL
WHERE t.deleted_at IS NULL
    AND t.transacted_at >= @start_date
    AND t.transacted_at <= @end_date
GROUP BY c.sub_id;

-- name: GetAccountBreakdown :many
SELECT
    a.id as account_id,
    a.name as account_name,
    t.transaction_type,
    COALESCE(SUM(t.enhanced_amount), 0) as total
FROM transactions t
LEFT JOIN accounts a ON t.account_id = a.id AND a.deleted_at IS NULL
WHERE t.deleted_at IS NULL
    AND t.transacted_at >= @start_date
    AND t.transacted_at <= @end_date
GROUP BY a.id, t.transaction_type;

-- name: GetTrends :many
SELECT
    DATE(transacted_at) as date,
    transaction_type,
    COALESCE(SUM(enhanced_amount), 0) as total
FROM transactions
WHERE deleted_at IS NULL
    AND transacted_at >= @start_date
    AND transacted_at <= @end_date
GROUP BY DATE(transacted_at), transaction_type;
