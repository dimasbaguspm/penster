-- name: CreateTransaction :one
INSERT INTO transactions (account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes)
VALUES (@account_id, @transfer_account_id, @category_id, @transaction_type, @title, @base_amount, @enhanced_amount, @currency, @currency_rate, @transacted_at, @notes)
RETURNING id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at;

-- name: GetTransactionByID :one
SELECT id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at
FROM transactions
WHERE id = @id AND deleted_at IS NULL;

-- name: GetTransactionBySubID :one
SELECT id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at
FROM transactions
WHERE sub_id = @sub_id AND deleted_at IS NULL;

-- name: ListTransactions :many
WITH base_query AS (
    SELECT id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at
    FROM transactions
    WHERE deleted_at IS NULL
),
search_query AS (
    SELECT * FROM base_query
    WHERE
        (@sub_id::uuid IS NULL OR sub_id = @sub_id)
        AND (@transaction_type::text IS NULL OR transaction_type = @transaction_type)
        AND (@q::text IS NULL OR title ILIKE '%' || @q || '%')
),
count_query AS (
    SELECT count(*) as total FROM search_query
),
paginated_query AS (
    SELECT
        id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at
    FROM search_query
    WHERE
        (@sort_by::text = 'title' AND (
            (@sort_order::text = 'asc' AND title > @cursor_title) OR
            (@sort_order::text = 'desc' AND title < @cursor_title) OR
            (@sort_order::text = 'asc' AND @cursor_title IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_title IS NULL)
        ))
        OR (@sort_by::text = 'transacted_at' AND (
            (@sort_order::text = 'asc' AND transacted_at > @cursor_transacted_at) OR
            (@sort_order::text = 'desc' AND transacted_at < @cursor_transacted_at) OR
            (@sort_order::text = 'asc' AND @cursor_transacted_at IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_transacted_at IS NULL)
        ))
        OR (@sort_by::text = 'created_at' AND (
            (@sort_order::text = 'asc' AND created_at > @cursor_created_at) OR
            (@sort_order::text = 'desc' AND created_at < @cursor_created_at) OR
            (@sort_order::text = 'asc' AND @cursor_created_at IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_created_at IS NULL)
        ))
        OR (@sort_by::text = 'amount' AND (
            (@sort_order::text = 'asc' AND enhanced_amount > @cursor_enhanced_amount) OR
            (@sort_order::text = 'desc' AND enhanced_amount < @cursor_enhanced_amount) OR
            (@sort_order::text = 'asc' AND @cursor_enhanced_amount IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_enhanced_amount IS NULL)
        ))
        OR @sort_by::text = ''
    ORDER BY
        CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN title END ASC,
        CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN title END DESC,
        CASE WHEN @sort_by::text = 'transacted_at' AND @sort_order::text = 'asc' THEN transacted_at END ASC,
        CASE WHEN @sort_by::text = 'transacted_at' AND @sort_order::text = 'desc' THEN transacted_at END DESC,
        CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at END ASC,
        CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at END DESC,
        CASE WHEN @sort_by::text = 'amount' AND @sort_order::text = 'asc' THEN enhanced_amount END ASC,
        CASE WHEN @sort_by::text = 'amount' AND @sort_order::text = 'desc' THEN enhanced_amount END DESC,
        CASE WHEN @sort_by::text = '' THEN id END ASC
    LIMIT NULLIF(@page_size, 0)
)
SELECT
    pq.id, pq.sub_id, pq.account_id, pq.transfer_account_id, pq.category_id, pq.transaction_type, pq.title, pq.base_amount, pq.enhanced_amount, pq.currency, pq.currency_rate, pq.transacted_at, pq.notes, pq.deleted_at, pq.created_at, pq.updated_at,
    cq.total
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
    transacted_at = COALESCE(@transacted_at, transacted_at),
    notes = @notes,
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at;

-- name: DeleteTransaction :one
UPDATE transactions
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, deleted_at, created_at, updated_at;
