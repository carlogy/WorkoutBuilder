-- name: CreateExerciseSets :exec
INSERT INTO exercise_sets (id, ordinal, workout_exerciseID, weight, reps, static_hold_time, created_at, modified_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NOW(),
    NOW()
);


-- name: GetExerciseSetsByExID :many
SELECT *
FROM
    exercise_sets
WHERE workout_exerciseID = $1
ORDER BY ordinal;


-- name: GetExerciseSetsByWorkoutID :many
Select
    es.id,
    es.ordinal,
    es.workout_exerciseID,
    es.weight,
    es.reps,
    es.static_hold_time,
    es.created_at,
    es.modified_at
FROM
    exercise_sets es
JOIN workout_exercise we
    ON es.workout_exerciseID = we.id
JOIN workout_blocks wb
    ON we.workout_blockID = wb.id
JOIN workouts w
    ON wb.workoutid = w.id
WHERE
    w.id = $1
ORDER BY
    es.ordinal;
