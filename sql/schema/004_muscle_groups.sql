-- +goose Up
CREATE TABLE muscle_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    body_part TEXT NOT NULL,
    muscle_group TEXT NOT NULL,
    muscle_name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE muscle_groups;
