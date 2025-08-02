-- name: CreateExercise :one
INSERT INTO exercises (id, name, exercise_type, equipment, primary_muscle_groups,secondary_muscle_groups, description, created_at, modified_at
)
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
    e.id = $1 LIMIT 1;

-- name: DeleteExerciseById :one
DELETE FROM exercises e
WHERE e.id = $1
RETURNING *;
