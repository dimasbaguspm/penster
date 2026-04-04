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
SELECT
    t.id, t.sub_id, t.account_id, t.transfer_account_id, t.category_id,
    t.transaction_type, t.title, t.base_amount, t.enhanced_amount,
    t.currency, t.currency_rate, t.transacted_at, t.notes,
    t.deleted_at, t.created_at, t.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id,
    cnt.total
FROM transactions t
LEFT JOIN accounts a ON t.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON t.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON t.category_id = c.id AND c.deleted_at IS NULL
CROSS JOIN (SELECT count(*) as total FROM transactions t2
    WHERE t2.deleted_at IS NULL
    AND ($1::uuid IS NULL OR t2.sub_id = $1)
    AND ($2 = 0 OR t2.account_id = $2)
    AND ($3 = 0 OR t2.category_id = $3)
    AND ($4 = '' OR t2.transaction_type = $4)
    AND ($5::text IS NULL OR t2.title ILIKE '%' || $5 || '%')
) cnt
WHERE t.deleted_at IS NULL
    AND ($1::uuid IS NULL OR t.sub_id = $1)
    AND ($2 = 0 OR t.account_id = $2)
    AND ($3 = 0 OR t.category_id = $3)
    AND ($4 = '' OR t.transaction_type = $4)
    AND ($5::text IS NULL OR t.title ILIKE '%' || $5 || '%')
    AND (
        ($6::text = 'title' AND (
            ($7::text = 'asc' AND t.title > $8) OR
            ($7::text = 'desc' AND t.title < $8) OR
            ($7::text = 'asc' AND $8 IS NULL) OR
            ($7::text = 'desc' AND $8 IS NULL)
        ))
        OR ($6::text = 'transacted_at' AND (
            ($7::text = 'asc' AND t.transacted_at > $9) OR
            ($7::text = 'desc' AND t.transacted_at < $9) OR
            ($7::text = 'asc' AND $9 IS NULL) OR
            ($7::text = 'desc' AND $9 IS NULL)
        ))
        OR ($6::text = 'created_at' AND (
            ($7::text = 'asc' AND t.created_at > $10) OR
            ($7::text = 'desc' AND t.created_at < $10) OR
            ($7::text = 'asc' AND $10 IS NULL) OR
            ($7::text = 'desc' AND $10 IS NULL)
        ))
        OR ($6::text = 'amount' AND (
            ($7::text = 'asc' AND t.enhanced_amount > $11) OR
            ($7::text = 'desc' AND t.enhanced_amount < $11) OR
            ($7::text = 'asc' AND $11 IS NULL) OR
            ($7::text = 'desc' AND $11 IS NULL)
        ))
        OR $6::text = ''
    )
ORDER BY
    CASE WHEN $6::text = 'title' AND $7::text = 'asc' THEN t.title END ASC,
    CASE WHEN $6::text = 'title' AND $7::text = 'desc' THEN t.title END DESC,
    CASE WHEN $6::text = 'transacted_at' AND $7::text = 'asc' THEN t.transacted_at END ASC,
    CASE WHEN $6::text = 'transacted_at' AND $7::text = 'desc' THEN t.transacted_at END DESC,
    CASE WHEN $6::text = 'created_at' AND $7::text = 'asc' THEN t.created_at END ASC,
    CASE WHEN $6::text = 'created_at' AND $7::text = 'desc' THEN t.created_at END DESC,
    CASE WHEN $6::text = 'amount' AND $7::text = 'asc' THEN t.enhanced_amount END ASC,
    CASE WHEN $6::text = 'amount' AND $7::text = 'desc' THEN t.enhanced_amount END DESC,
    CASE WHEN $6::text = '' THEN t.id END ASC
LIMIT NULLIF($12, 0);

-- name: UpdateTransaction :one
UPDATE transactions
SET
    account_id = COALESCE(NULLIF(@account_id, 0), account_id),
    transfer_account_id = COALESCE(@transfer_account_id, transfer_account_id),
    category_id = COALESCE(NULLIF(@category_id, 0), category_id),
    transaction_type = COALESCE(NULLIF(@transaction_type, ''), transaction_type),
    title = COALESCE(NULLIF(@title, ''), title),
    base_amount = COALESCE(NULLIF(@base_amount, 0), base_amount),
    enhanced_amount = COALESCE(@enhanced_amount, enhanced_amount),
    currency = COALESCE(NULLIF(@currency, ''), currency),
    currency_rate = COALESCE(@currency_rate, currency_rate),
    notes = COALESCE(@notes, notes),
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id;

-- name: DeleteTransaction :one
UPDATE transactions
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id;
