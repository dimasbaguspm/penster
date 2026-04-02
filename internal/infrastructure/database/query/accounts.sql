-- name: CreateAccount :one
INSERT INTO accounts (sub_id, name, type, balance)
VALUES (gen_random_uuid(), @name, @type, @balance)
RETURNING id, sub_id, name, type, balance, deleted_at, created_at, updated_at;

-- name: GetAccountByID :one
SELECT id, sub_id, name, type, balance, deleted_at, created_at, updated_at
FROM accounts
WHERE id = @id AND deleted_at IS NULL;

-- name: GetAccountBySubID :one
SELECT id, sub_id, name, type, balance, deleted_at, created_at, updated_at
FROM accounts
WHERE sub_id = @sub_id AND deleted_at IS NULL;

-- name: ListAccounts :many
WITH base_query AS (
    SELECT id, sub_id, name, type, balance, deleted_at, created_at, updated_at
    FROM accounts
    WHERE deleted_at IS NULL
),
search_query AS (
    SELECT * FROM base_query
    WHERE
        (@sub_id::uuid IS NULL OR sub_id = @sub_id)
        AND (@q::text IS NULL OR name ILIKE '%' || @q || '%')
),
count_query AS (
    SELECT count(*) as total FROM search_query
),
paginated_query AS (
    SELECT
        id, sub_id, name, type, balance, deleted_at, created_at, updated_at
    FROM search_query
    WHERE
        (@sort_by::text = 'name' AND (
            (@sort_order::text = 'asc' AND name > @cursor_name) OR
            (@sort_order::text = 'desc' AND name < @cursor_name) OR
            (@sort_order::text = 'asc' AND @cursor_name IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_name IS NULL)
        ))
        OR (@sort_by::text = 'created_at' AND (
            (@sort_order::text = 'asc' AND created_at > @cursor_created_at) OR
            (@sort_order::text = 'desc' AND created_at < @cursor_created_at) OR
            (@sort_order::text = 'asc' AND @cursor_created_at IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_created_at IS NULL)
        ))
        OR (@sort_by::text = 'updated_at' AND (
            (@sort_order::text = 'asc' AND updated_at > @cursor_updated_at) OR
            (@sort_order::text = 'desc' AND updated_at < @cursor_updated_at) OR
            (@sort_order::text = 'asc' AND @cursor_updated_at IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_updated_at IS NULL)
        ))
        OR (@sort_by::text = 'balance' AND (
            (@sort_order::text = 'asc' AND balance > @cursor_balance) OR
            (@sort_order::text = 'desc' AND balance < @cursor_balance) OR
            (@sort_order::text = 'asc' AND @cursor_balance IS NULL) OR
            (@sort_order::text = 'desc' AND @cursor_balance IS NULL)
        ))
        OR @sort_by::text = ''
    ORDER BY
        CASE WHEN @sort_by::text = 'name' AND @sort_order::text = 'asc' THEN name END ASC,
        CASE WHEN @sort_by::text = 'name' AND @sort_order::text = 'desc' THEN name END DESC,
        CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at END ASC,
        CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at END DESC,
        CASE WHEN @sort_by::text = 'updated_at' AND @sort_order::text = 'asc' THEN updated_at END ASC,
        CASE WHEN @sort_by::text = 'updated_at' AND @sort_order::text = 'desc' THEN updated_at END DESC,
        CASE WHEN @sort_by::text = 'balance' AND @sort_order::text = 'asc' THEN balance END ASC,
        CASE WHEN @sort_by::text = 'balance' AND @sort_order::text = 'desc' THEN balance END DESC,
        CASE WHEN @sort_by::text = '' THEN id END ASC
    LIMIT NULLIF(@page_size, 0)
)
SELECT
    pq.id, pq.sub_id, pq.name, pq.type, pq.balance, pq.deleted_at, pq.created_at, pq.updated_at,
    cq.total
FROM paginated_query pq, count_query cq;

-- name: UpdateAccount :one
UPDATE accounts
SET
    name = COALESCE(@name, name),
    type = COALESCE(@type, type),
    balance = COALESCE(@balance, balance),
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, name, type, balance, deleted_at, created_at, updated_at;

-- name: DeleteAccount :one
UPDATE accounts
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, name, type, balance, deleted_at, created_at, updated_at;
