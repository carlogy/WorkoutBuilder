-- name: CreateExercise :one
INSERT INTO exercises (id, name, exercise_type, equipment, description, has_primary_muscles, has_secondary_muscles, created_at, modified_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetExercises :many
SELECT
    *
FROM
    exercises e;

-- name: GetExerciseById :one
SELECT
    *
FROM
    exercises e
WHERE
    e.id = $1
LIMIT 1;

-- name: DeleteExerciseById :one
DELETE FROM exercises e
WHERE e.id = $1
RETURNING *;

-- name: CheckExerciseExists :one
SELECT EXISTS (
SELECT 1
FROM
    exercises
WHERE name = $1
);
