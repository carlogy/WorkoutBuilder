-- name: CreatePlan :one
INSERT INTO plans (id, name, goal, days, duration, description, experience_level, created_at, modified_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NOW(),
    NOW()
)
RETURNING *;
