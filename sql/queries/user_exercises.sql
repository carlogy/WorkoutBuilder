-- name: CreateUserExercise :one
INSERT INTO user_exercises(id, userId, exerciseId, sets_weight, rest, duration, decline_incline, notes, created_at, modified_at)
Values(
    gen_random_uuid(),
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


-- name: GetUserExerciseRecordById :one
SELECT
    *
FROM
    user_exercises ue
WHERE
    ue.id = $1 LIMIT 1;

-- name: UpdateUserExcerciseRecordById :one
UPDATE
    user_exercises
SET
    sets_weight = $1,
    rest = $2,
    duration = $3,
    decline_incline = $4,
    notes = $5,
    modified_at = NOW()
WHERE
    id = $6
RETURNING *;

-- name: DeleteUserExerciseRecordById :one
DELETE FROM user_exercises ue
WHERE ue.id = $1
RETURNING *;
