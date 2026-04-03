-- name: CreateCategory :one
INSERT INTO categories (sub_id, name, type)
VALUES (gen_random_uuid(), @name, @type)
RETURNING id, sub_id, name, type, deleted_at, created_at, updated_at;

-- name: GetCategoryByID :one
SELECT id, sub_id, name, type, deleted_at, created_at, updated_at
FROM categories
WHERE id = @id AND deleted_at IS NULL;

-- name: GetCategoryBySubID :one
SELECT id, sub_id, name, type, deleted_at, created_at, updated_at
FROM categories
WHERE sub_id = @sub_id AND deleted_at IS NULL;

-- name: ListCategories :many
SELECT
    c.id, c.sub_id, c.name, c.type, c.deleted_at, c.created_at, c.updated_at,
    cnt.total
FROM categories c
CROSS JOIN (SELECT count(*) as total FROM categories c2
    WHERE c2.deleted_at IS NULL
    AND ($1::uuid IS NULL OR c2.sub_id = $1::uuid)
    AND ($2::text IS NULL OR c2.name ILIKE '%' || $2::text || '%')
) cnt
WHERE c.deleted_at IS NULL
    AND ($1::uuid IS NULL OR c.sub_id = $1::uuid)
    AND ($2::text IS NULL OR c.name ILIKE '%' || $2::text || '%')
    AND (
        ($3::text = 'name' AND (
            ($4::text = 'asc' AND c.name > $5::text) OR
            ($4::text = 'desc' AND c.name < $5::text) OR
            ($4::text = 'asc' AND $5::text IS NULL) OR
            ($4::text = 'desc' AND $5::text IS NULL)
        ))
        OR ($3::text = 'created_at' AND (
            ($4::text = 'asc' AND c.created_at > $6::timestamptz) OR
            ($4::text = 'desc' AND c.created_at < $6::timestamptz) OR
            ($4::text = 'asc' AND $6::timestamptz IS NULL) OR
            ($4::text = 'desc' AND $6::timestamptz IS NULL)
        ))
        OR ($3::text = 'updated_at' AND (
            ($4::text = 'asc' AND c.updated_at > $7::timestamptz) OR
            ($4::text = 'desc' AND c.updated_at < $7::timestamptz) OR
            ($4::text = 'asc' AND $7::timestamptz IS NULL) OR
            ($4::text = 'desc' AND $7::timestamptz IS NULL)
        ))
        OR $3::text = ''
    )
ORDER BY
    CASE WHEN $3::text = 'name' AND $4::text = 'asc' THEN c.name END ASC,
    CASE WHEN $3::text = 'name' AND $4::text = 'desc' THEN c.name END DESC,
    CASE WHEN $3::text = 'created_at' AND $4::text = 'asc' THEN c.created_at END ASC,
    CASE WHEN $3::text = 'created_at' AND $4::text = 'desc' THEN c.created_at END DESC,
    CASE WHEN $3::text = 'updated_at' AND $4::text = 'asc' THEN c.updated_at END ASC,
    CASE WHEN $3::text = 'updated_at' AND $4::text = 'desc' THEN c.updated_at END DESC,
    CASE WHEN $3::text = '' THEN c.id END ASC
LIMIT NULLIF($8, 0);

-- name: UpdateCategory :one
UPDATE categories
SET
    name = COALESCE(NULLIF(@name, ''), name),
    type = COALESCE(NULLIF(@type, ''), type),
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, name, type, deleted_at, created_at, updated_at;

-- name: DeleteCategory :one
UPDATE categories
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING id, sub_id, name, type, deleted_at, created_at, updated_at;
