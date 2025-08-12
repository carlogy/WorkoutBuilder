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

-- name: GetWorkouts :many
SELECT *
FROM
    workouts;

-- name: GetWorkoutByID :one
SELECT
    *
FROM
    workouts w
WHERE
    w.id = $1;

-- name: DeleteWorkoutByID :one
DELETE FROM workouts
WHERE id = $1
RETURNING *;
