-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    exercise_type TEXT NOT NULL,
    equipment TEXT NOT NULL,
    primary_muscle_groups JSONB,
    secondary_muscle_groups JSONB,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
Drop TABLE exercises;