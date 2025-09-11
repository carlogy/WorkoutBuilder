package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type ExerciseRepository struct {
	dbQ *db.Queries
	db  *sql.DB
}

type ExercisesRepository interface {
	CreateDBExercise(ctx context.Context, exParams db.CreateExerciseParams) (db.Exercise, error)
	GetExerciseByID(ctx context.Context, id id.UUID) (db.Exercise, error)
	GetMuscleGroupByName(ctx context.Context, name string) db.MuscleGroup
	CheckExerciseExists(ctx context.Context, name string) (bool, error)
	GetMuscleGroupsByExerciseID(ctx context.Context, exId id.UUID) ([]db.GetMuscleGroupsByExerciseIDRow, error)
}

func NewExerciseRepository(dbq *db.Queries, db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{
		dbQ: dbq,
		db:  db,
	}
}

func (er *ExerciseRepository) CreateDBExercise(ctx context.Context, exParams db.CreateExerciseParams, muscleGroupSlice []db.CreateMuscleGroupParams, exMusGroupSlice []db.CreateExerciseMuscleGroupsParams) (db.Exercise, error) {

	tx, err := er.db.BeginTx(ctx, nil)
	if err != nil {
		return db.Exercise{}, err
	}
	defer func() {
		if err != nil {
			commitErr := tx.Rollback()
			if commitErr != nil {
				fmt.Println("Error rolling back: ", commitErr)
			}
		}
	}()
	qtx := er.dbQ.WithTx(tx)

	dbEX, err := qtx.CreateExercise(ctx, exParams)
	if err != nil {
		return db.Exercise{}, err
	}

	for _, muscleGroup := range muscleGroupSlice {
		_, err := qtx.CreateMuscleGroup(ctx, muscleGroup)
		if err != nil {
			return db.Exercise{}, err
		}
	}

	for _, exmuscle := range exMusGroupSlice {
		err := qtx.CreateExerciseMuscleGroups(ctx, exmuscle)
		if err != nil {
			return db.Exercise{}, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return db.Exercise{}, err
	}
	return dbEX, nil
}

func (er *ExerciseRepository) GetExerciseByID(ctx context.Context, id id.UUID) (db.Exercise, error) {
	return er.dbQ.GetExerciseById(ctx, id)
}

func (er *ExerciseRepository) GetMuscleGroupByName(ctx context.Context, name string) (db.MuscleGroup, error) {

	mg, err := er.dbQ.GetMuscleGroupByMuscleName(ctx, name)
	if err != nil {
		return db.MuscleGroup{}, err
	}
	return mg, nil
}

func (er *ExerciseRepository) CheckExerciseExists(ctx context.Context, name string) (bool, error) {
	return er.dbQ.CheckExerciseExists(ctx, name)
}

func (er *ExerciseRepository) GetMuscleGroupsByExerciseID(ctx context.Context, exId id.UUID) ([]db.GetMuscleGroupsByExerciseIDRow, error) {
	return er.dbQ.GetMuscleGroupsByExerciseID(ctx, exId)
}
