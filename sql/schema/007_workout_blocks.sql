-- +goose Up
CREATE TABLE workout_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workoutID UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    restSeconds_after_block INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE workout_blocks;
