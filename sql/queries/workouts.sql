-- name: CreateWorkOut :one
INSERT INTO workouts (id, name, description, exercises, created_at, modified_at)
Values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING *;
