-- name: CreateMuscleGroup :one
INSERT INTO muscle_groups (id, body_part, muscle_group, muscle_name, created_at, modified_at)
VALUES (
$1,
$2,
$3,
$4,
NOW(),
NOW()
)
RETURNING *;

-- name: GetMuscleGroupByMuscleName :one
SELECT *
FROM
    muscle_groups mg
WHERE
    mg.muscle_name = $1;

-- name: GetMuscleGroupsByExerciseID :many
Select
     sqlc.embed(mg), sqlc.embed(emg)
FROM
    muscle_groups mg
JOIN exercise_muscle_groups emg
    ON emg.muscle_groups_id = mg.id
WHERE
    emg.exercise_id = $1
ORDER BY emg.primary_muscle, emg.secondary_muscle;

-- name: GetMuscleGroupsForAllExercises :many
Select
     sqlc.embed(mg), sqlc.embed(emg)
FROM
    muscle_groups mg
JOIN exercise_muscle_groups emg
    ON emg.muscle_groups_id = mg.id
ORDER BY emg.primary_muscle, emg.secondary_muscle;
