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
