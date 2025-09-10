-- +goose UP
CREATE TABLE exercise_muscle_groups (
    id SERIAL PRIMARY KEY,
    exercise_ID UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    muscle_groups_ID UUID NOT NULL REFERENCES muscle_groups(id),
    primary_muscle boolean,
    secondary_muscle boolean,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE exercise_muscle_Groups;
