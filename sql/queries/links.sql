-- name: CreateNewLink :one
INSERT INTO links (id, created_at, updated_at, code, link_url, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetNewLink :one
SELECT code FROM links WHERE user_id = $1 ORDER BY created_at DESC
LIMIT 1;

-- name: GetLink :one
SELECT * FROM links WHERE code = $1;

