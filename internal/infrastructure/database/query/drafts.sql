-- name: CreateDraft :one
INSERT INTO drafts (account_id, transfer_account_id, category_id, transaction_type, title, base_amount, enhanced_amount, currency, currency_rate, transacted_at, notes, source, status)
VALUES (@account_id, @transfer_account_id, @category_id, @transaction_type, @title, @base_amount, @enhanced_amount, @currency, @currency_rate, @transacted_at, @notes, @source, @status)
RETURNING id;

-- name: GetDraftByID :one
SELECT
    d.id, d.sub_id, d.account_id, d.transfer_account_id, d.category_id,
    d.transaction_type, d.title, d.base_amount, d.enhanced_amount,
    d.currency, d.currency_rate, d.transacted_at, d.notes,
    d.source, d.status, d.confirmed_at, d.rejected_at,
    d.deleted_at, d.created_at, d.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id
FROM drafts d
LEFT JOIN accounts a ON d.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON d.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON d.category_id = c.id AND c.deleted_at IS NULL
WHERE d.id = @id AND d.deleted_at IS NULL;

-- name: GetDraftBySubID :one
SELECT
    d.id, d.sub_id, d.account_id, d.transfer_account_id, d.category_id,
    d.transaction_type, d.title, d.base_amount, d.enhanced_amount,
    d.currency, d.currency_rate, d.transacted_at, d.notes,
    d.source, d.status, d.confirmed_at, d.rejected_at,
    d.deleted_at, d.created_at, d.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id
FROM drafts d
LEFT JOIN accounts a ON d.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON d.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON d.category_id = c.id AND c.deleted_at IS NULL
WHERE d.sub_id = @sub_id AND d.deleted_at IS NULL;

-- name: ListDrafts :many
SELECT
    d.id, d.sub_id, d.account_id, d.transfer_account_id, d.category_id,
    d.transaction_type, d.title, d.base_amount, d.enhanced_amount,
    d.currency, d.currency_rate, d.transacted_at, d.notes,
    d.source, d.status, d.confirmed_at, d.rejected_at,
    d.deleted_at, d.created_at, d.updated_at,
    a.sub_id as account_sub_id,
    ta.sub_id as transfer_account_sub_id,
    c.sub_id as category_sub_id,
    cnt.total
FROM drafts d
LEFT JOIN accounts a ON d.account_id = a.id AND a.deleted_at IS NULL
LEFT JOIN accounts ta ON d.transfer_account_id = ta.id AND ta.deleted_at IS NULL
LEFT JOIN categories c ON d.category_id = c.id AND c.deleted_at IS NULL
CROSS JOIN (SELECT count(*) as total FROM drafts d2
    WHERE d2.deleted_at IS NULL
    AND (''::text = $1 OR d2.source = $1)
    AND (''::text = $2 OR d2.status = $2)
) cnt
WHERE d.deleted_at IS NULL
    AND (''::text = $1 OR d.source = $1)
    AND (''::text = $2 OR d.status = $2)
ORDER BY d.created_at DESC
LIMIT NULLIF(@page_size, 0);

-- name: UpdateDraft :one
UPDATE drafts
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
    transacted_at = COALESCE(@transacted_at, transacted_at),
    notes = COALESCE(@notes, notes),
    updated_at = NOW()
WHERE sub_id = @sub_id AND deleted_at IS NULL
RETURNING id;

-- name: UpdateDraftStatus :one
WITH input AS (
    SELECT
        @status::VARCHAR as new_status,
        CASE WHEN @status = 'confirmed' THEN NOW() ELSE NULL END as new_confirmed,
        CASE WHEN @status = 'rejected' THEN NOW() ELSE NULL END as new_rejected
)
UPDATE drafts d
SET
    status = i.new_status,
    confirmed_at = COALESCE(i.new_confirmed, d.confirmed_at),
    rejected_at = COALESCE(i.new_rejected, d.rejected_at),
    updated_at = NOW()
FROM input i
WHERE d.sub_id = @sub_id AND d.deleted_at IS NULL
RETURNING d.id;

-- name: SoftDeleteDraft :one
UPDATE drafts
SET deleted_at = NOW(), updated_at = NOW()
WHERE sub_id = @sub_id AND deleted_at IS NULL
RETURNING id;
