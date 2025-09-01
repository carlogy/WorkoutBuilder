-- +goose Up
CREATE TABLE workout_exercise (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_blockID UUID NOT NULL REFERENCES workout_blocks(id) ON DELETE CASCADE,
    exerciseID UUID NOT NULL REFERENCES exercises(id),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE workout_exercise;
