-- name: CreateExerciseMuscleGroups :exec
INSERT INTO exercise_muscle_groups (exercise_ID, muscle_groups_ID, primary_muscle, secondary_muscle, created_at, modified_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
);
