-- #nosec G101 -- This is a false positive
-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() +  INTERVAL '60 days'
)
RETURNING *;

-- #nosec G101 -- This is a false positive
-- name: GetRefreshToken :one
SELECT *
    FROM users u
INNER JOIN
    refresh_tokens rt
ON
    rt.user_id = u.user_id
WHERE
    rt.token = $1
AND
    rt.revoked_at is not null
AND
    rt.expires_at > NOW();

-- #nosec G101 -- This is a false positive
-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token = $1
WHERE token = $2
RETURNING *;
