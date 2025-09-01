-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    exercise_type TEXT NOT NULL,
    equipment TEXT NOT NULL,
    description TEXT,
    has_primary_muscles boolean,
    has_secondary_muscles boolean,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
Drop TABLE exercises;
