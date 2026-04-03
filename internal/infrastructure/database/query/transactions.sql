-- name: CreateTransaction :one
INSERT INTO transactions (account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes)
VALUES (@account_id, @transfer_account_id, @category_id, @transaction_type, @title, @base_amount, @enhanced_amount, @currency, @currency_rate, @transacted_at, @notes)
RETURNING id;

-- name: GetTransactionByID :one
SELECT
    t.id, t.sub_id, t.account_id, t.transfer_account_id, t.category_id,
    t.transaction_type, t.title, t.base_amount, t.enhanced_amount,
    t.currency, t.currency_rate, t.transacted_at, t.notes,
    t.deleted_at, t.created_at, t.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id
FROM transactions t
LEFT JOIN accounts a ON t.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON t.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON t.category_id = c.id AND c.deleted_at IS NULL
WHERE t.id = @id AND t.deleted_at IS NULL;

-- name: GetTransactionBySubID :one
SELECT
    t.id, t.sub_id, t.account_id, t.transfer_account_id, t.category_id,
    t.transaction_type, t.title, t.base_amount, t.enhanced_amount,
    t.currency, t.currency_rate, t.transacted_at, t.notes,
    t.deleted_at, t.created_at, t.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id
FROM transactions t
LEFT JOIN accounts a ON t.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON t.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON t.category_id = c.id AND c.deleted_at IS NULL
WHERE t.sub_id = @sub_id AND t.deleted_at IS NULL;

-- name: ListTransactions :many
WITH base_query AS (
    SELECT t.id, t.sub_id, t.account_id, t.transfer_account_id, t.category_id,
        t.transaction_type, t.title, t.base_amount, t.enhanced_amount,
        t.currency, t.currency_rate, t.transacted_at, t.notes,
        t.deleted_at, t.created_at, t.updated_at,
        a.sub_id as account_sub_id,
        ta.sub_id as transfer_account_sub_id,
        c.sub_id as category_sub_id
    FROM transactions t
    LEFT JOIN accounts a ON t.account_id = a.id AND a.deleted_at IS NULL
    LEFT JOIN accounts ta ON t.transfer_account_id = ta.id AND ta.deleted_at IS NULL
    LEFT JOIN categories c ON t.category_id = c.id AND c.deleted_at IS NULL
    WHERE t.deleted_at IS NULL
),
search_query AS (
    SELECT * FROM base_query
    WHERE
        ($1::uuid IS NULL OR sub_id = $1)
        AND ($2::int IS NULL OR account_id = $2)
        AND ($3::int IS NULL OR category_id = $3)
        AND ($4::text IS NULL OR transaction_type = $4)
        AND ($5::text IS NULL OR title ILIKE '%' || $5 || '%')
),
count_query AS (
    SELECT count(*) as total FROM search_query
),
paginated_query AS (
    SELECT *
    FROM search_query
    WHERE
        ($6::text = 'title' AND (
            ($7::text = 'asc' AND title > $8) OR
            ($7::text = 'desc' AND title < $8) OR
            ($7::text = 'asc' AND $8 IS NULL) OR
            ($7::text = 'desc' AND $8 IS NULL)
        ))
        OR ($6::text = 'transacted_at' AND (
            ($7::text = 'asc' AND transacted_at > $9) OR
            ($7::text = 'desc' AND transacted_at < $9) OR
            ($7::text = 'asc' AND $9 IS NULL) OR
            ($7::text = 'desc' AND $9 IS NULL)
        ))
        OR ($6::text = 'created_at' AND (
            ($7::text = 'asc' AND created_at > $10) OR
            ($7::text = 'desc' AND created_at < $10) OR
            ($7::text = 'asc' AND $10 IS NULL) OR
            ($7::text = 'desc' AND $10 IS NULL)
        ))
        OR ($6::text = 'amount' AND (
            ($7::text = 'asc' AND enhanced_amount > $11) OR
            ($7::text = 'desc' AND enhanced_amount < $11) OR
            ($7::text = 'asc' AND $11 IS NULL) OR
            ($7::text = 'desc' AND $11 IS NULL)
        ))
        OR $4::text = ''
    ORDER BY
        CASE WHEN $4::text = 'title' AND $5::text = 'asc' THEN title END ASC,
        CASE WHEN $4::text = 'title' AND $5::text = 'desc' THEN title END DESC,
        CASE WHEN $4::text = 'transacted_at' AND $5::text = 'asc' THEN transacted_at END ASC,
        CASE WHEN $4::text = 'transacted_at' AND $5::text = 'desc' THEN transacted_at END DESC,
        CASE WHEN $4::text = 'created_at' AND $5::text = 'asc' THEN created_at END ASC,
        CASE WHEN $4::text = 'created_at' AND $5::text = 'desc' THEN created_at END DESC,
        CASE WHEN $4::text = 'amount' AND $5::text = 'asc' THEN enhanced_amount END ASC,
        CASE WHEN $4::text = 'amount' AND $5::text = 'desc' THEN enhanced_amount END DESC,
        CASE WHEN $4::text = '' THEN id END ASC
    LIMIT NULLIF($12, 0)
)
SELECT pq.*, cq.total
FROM paginated_query pq, count_query cq;

-- name: UpdateTransaction :one
UPDATE transactions
SET
    account_id = COALESCE(@account_id, account_id),
    transfer_account_id = @transfer_account_id,
    category_id = @category_id,
    transaction_type = COALESCE(@transaction_type, transaction_type),
    title = COALESCE(@title, title),
    base_amount = COALESCE(@base_amount, base_amount),
    enhanced_amount = @enhanced_amount,
    currency = COALESCE(@currency, currency),
    currency_rate = COALESCE(@currency_rate, currency_rate),
    notes = @notes,
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id;

-- name: DeleteTransaction :one
UPDATE transactions
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id;
