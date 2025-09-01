-- +goose Up
CREATE TABLE exercise_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exerciseID UUID NOT NULL REFERENCES workout_exercise(id),
    weight DECIMAL,
    reps INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE exercise_sets;
