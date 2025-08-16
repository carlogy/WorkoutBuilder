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


-- name: UpdateRefreshToken :one
-- #nosec G101 -- This is a false positive
UPDATE refresh_tokens
SET token = $1
WHERE token = $2
    AND
        revoked_at is null
    AND
        expires_at > NOW()
RETURNING *;


-- name: RevokeRefreshToken :one
-- #nosec G101 -- This is a false positive
UPDATE refresh_tokens
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE
    token = $1
RETURNING *;
