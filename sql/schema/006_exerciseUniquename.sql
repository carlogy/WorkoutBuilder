-- +goose Up
ALTER TABLE exercises
ADD CONSTRAINT unique_exercise_name
UNIQUE (name);

-- +goose Down
ALTER TABLE exercises
DROP CONSTRAINT unique_exercise_name;
