-- name: CreateWorkoutExercises :exec
INSERT INTO workout_exercise (id, ordinal, workout_blockID, exerciseID, notes, created_at, modified_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetWorkoutExercisesByWorkoutID :many
SELECT
    we.id,
    we.ordinal,
    we.workout_blockID,
    we.exerciseID,
    we.notes,
    we.created_at,
    we.modified_at,
    wb.workoutid AS workoutID
FROM
    workout_exercise we
JOIN workout_blocks wb
    ON we.workout_blockID = wb.id
JOIN workouts w
    on wb.workoutid = w.id
WHERE
    w.id = $1
ORDER BY
    we.ordinal;
