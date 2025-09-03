package repositories

import (
	"context"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type WorkoutRepository struct {
	db *db.Queries
}

type WorkOutRepository interface {
	CreateDBWO(ctx context.Context, workout *db.CreateWorkOutParams) (db.Workout, error)
	GetWODBById(ctx context.Context, id id.UUID) (db.Workout, error)
	DeleteDBWOById(ctx context.Context, id id.UUID) (db.Workout, error)
}

func NewWorkoutRepository(db *db.Queries) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (wo *WorkoutRepository) CreateDBWO(ctx context.Context, workout *db.CreateWorkOutParams) (db.Workout, error) {

	dbWO, err := wo.db.CreateWorkOut(ctx, *workout)
	if err != nil {
		return db.Workout{}, err
	}

	return dbWO, nil
}

func (wo *WorkoutRepository) GetWODBById(ctx context.Context, id id.UUID) (db.Workout, error) {
	return db.Workout{}, nil
}
func (wo *WorkoutRepository) DeleteDBWOById(ctx context.Context, id id.UUID) (db.Workout, error) {
	return db.Workout{}, nil
}
