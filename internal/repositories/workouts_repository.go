package repositories

import (
	"context"
	"database/sql"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type WorkoutRepository struct {
	dbQ *db.Queries
	db  *sql.DB
}

type WorkOutRepository interface {
	CreateDBWO(ctx context.Context, workout *db.CreateWorkOutParams, woBlockSlice []db.CreateWorkoutBlocksParams, woExSlice []db.CreateWorkoutExercisesParams, woExSets []db.CreateExerciseSetsParams) (db.Workout, error)
	GetWODBById(ctx context.Context, id id.UUID) ([]db.GetWorkoutByIDRow, error)
	DeleteDBWOById(ctx context.Context, id id.UUID) (db.Workout, error)
	GetDBWOExByWOID(ctx context.Context, workoutID id.UUID) ([]db.GetWorkoutExercisesByWorkoutIDRow, error)
	GetDBWOBlocksByWOID(ctx context.Context, woID id.UUID) ([]db.WorkoutBlock, error)
	GetExSetsByWOID(ctx context.Context, exID id.UUID) ([]db.GetExerciseSetsByWorkoutIDRow, error)
}

func NewWorkoutRepository(dbq *db.Queries, db *sql.DB) *WorkoutRepository {
	return &WorkoutRepository{dbQ: dbq, db: db}
}

func (wo *WorkoutRepository) CreateDBWO(ctx context.Context, workout *db.CreateWorkOutParams, woBlockSlice []db.CreateWorkoutBlocksParams, woExSlice []db.CreateWorkoutExercisesParams, woExSets []db.CreateExerciseSetsParams) (db.Workout, error) {

	tx, err := wo.db.Begin()
	if err != nil {
		return db.Workout{}, err
	}
	defer tx.Rollback()
	qtx := wo.dbQ.WithTx(tx)

	dbWO, err := qtx.CreateWorkOut(ctx, *workout)
	if err != nil {
		return db.Workout{}, err
	}

	for _, woBlock := range woBlockSlice {
		if err := qtx.CreateWorkoutBlocks(ctx, woBlock); err != nil {
			return db.Workout{}, err
		}
	}

	for _, ex := range woExSlice {

		if err := qtx.CreateWorkoutExercises(ctx, ex); err != nil {
			return db.Workout{}, err
		}
	}

	for _, exSet := range woExSets {

		if err := qtx.CreateExerciseSets(ctx, exSet); err != nil {
			return db.Workout{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return db.Workout{}, err
	}
	return dbWO, nil
}

func (wo *WorkoutRepository) GetWODBById(ctx context.Context, id id.UUID) ([]db.GetWorkoutByIDRow, error) {

	dbWO, err := wo.dbQ.GetWorkoutByID(ctx, id)
	if err != nil {
		return []db.GetWorkoutByIDRow{}, err
	}
	return dbWO, nil
}

func (wo *WorkoutRepository) DeleteDBWOById(ctx context.Context, id id.UUID) (db.Workout, error) {
	return db.Workout{}, nil
}

func (wo *WorkoutRepository) GetDBWOExByWOID(ctx context.Context, workoutID id.UUID) ([]db.GetWorkoutExercisesByWorkoutIDRow, error) {
	woExSlice, err := wo.dbQ.GetWorkoutExercisesByWorkoutID(ctx, workoutID)
	if err != nil {
		return []db.GetWorkoutExercisesByWorkoutIDRow{}, err
	}

	return woExSlice, nil
}

func (wo *WorkoutRepository) GetDBWOBlocksByWOID(ctx context.Context, woID id.UUID) ([]db.WorkoutBlock, error) {

	woBlocks, err := wo.dbQ.GetWorkoutBlocksByWOID(ctx, woID)
	if err != nil {
		return []db.WorkoutBlock{}, err
	}
	return woBlocks, nil
}

func (wo *WorkoutRepository) GetExSetsByWOID(ctx context.Context, workoutID id.UUID) ([]db.GetExerciseSetsByWorkoutIDRow, error) {

	exSets, err := wo.dbQ.GetExerciseSetsByWorkoutID(ctx, workoutID)
	if err != nil {
		return []db.GetExerciseSetsByWorkoutIDRow{}, err
	}
	return exSets, nil
}
