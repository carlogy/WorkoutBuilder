-- name: CreateWorkOut :one
INSERT INTO workouts (id, name, description,created_at, modified_at)
Values(
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

-- name: GetWorkoutByID :many
SELECT
    sqlc.embed(workouts),
    sqlc.embed(workout_blocks), sqlc.embed(workout_exercise), sqlc.embed(exercise_sets)
FROM workouts
JOIN workout_blocks on workout_blocks.workoutid = workouts.id
JOIN workout_exercise on workout_exercise.workout_blockid = workout_blocks.id
JOIN exercise_sets  on exercise_sets.workout_exerciseid = workout_exercise.id
WHERE workouts.id = $1;

-- name: DeleteWorkoutByID :one
DELETE FROM workouts
WHERE id = $1
RETURNING *;
