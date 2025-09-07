-- name: CreateWorkoutBlocks :exec
INSERT INTO workout_blocks (id, ordinal, workoutID, restSeconds_after_block, created_at, modified_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetWorkoutBlocksByWOID :many
SELECT *
FROM
    workout_blocks
WHERE
    workoutID = $1
ORDER BY ordinal;
