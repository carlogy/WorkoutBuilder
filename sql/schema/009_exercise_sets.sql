-- +goose Up
CREATE TABLE exercise_sets (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    workout_exerciseID UUID NOT NULL,
    ordinal INTEGER NOT NULL,
    weight DECIMAL(10,2),
    reps INTEGER,
    static_hold_time INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE exercise_sets;
