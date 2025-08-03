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
