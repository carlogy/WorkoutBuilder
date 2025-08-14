-- name: StoreRefreshToken :one
-- #nosec G101 -- This is a false positive
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() +  INTERVAL '60 days'
)
RETURNING *;
