-- +goose Up
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    Name TEXT NOT NULL,
    Description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);


-- +goose Down
DROP TABLE workouts;
