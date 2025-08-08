-- +goose Up
CREATE TABLE user_exercises (
    id UUID PRIMARY KEY,
    userId UUID NOT NULL REFERENCES users(id),
    exerciseId UUID NOT NULL REFERENCES exercises(id),
    sets_weight JSONB,
    rest BIGINT,
    duration BIGINT,
    decline_incline BIGINT,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE user_exercises;
