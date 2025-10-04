-- +goose Up
CREATE TABLE workouts (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    Name TEXT NOT NULL,
    Description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);


-- +goose Down
DROP TABLE workouts;
