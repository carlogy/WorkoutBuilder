-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, password, created_at, modified_at)
VALUES(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
RETURNING id, first_name, last_name, email, created_at, modified_at;


-- name: GetUserByEmail :one
SELECT
   *
FROM
    users u
WHERE
    u.email = $1;

-- name: UpdateUserById :one
UPDATE
    users
SET
    first_name = $1,
    last_name = $2,
    email = $3,
    password = $4,
    modified_at = NOW()
WHERE
    id = $5
RETURNING id, first_name, last_name, email, created_at, modified_at;

-- name: DeleteUserById :one
DELETE FROM
    users u
WHERE
    u.id = $1
RETURNING id, first_name, last_name, email, created_at, modified_at;
