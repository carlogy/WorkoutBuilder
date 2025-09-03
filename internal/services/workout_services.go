package services

import (
	"context"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
	id "github.com/google/uuid"
)

type WorkoutService struct {
	workoutRepo repositories.WorkOutRepository
	Secret      string
}

type Set struct {
	Weight int `json:"weight"`
	Reps   int `json:"reps"`
}

type WorkoutExercise struct {
	ExerciseID id.UUID `json:"exerciseID"`
	Sets       []Set   `json:"sets"`
	Notes      string  `JSON:"notes,omitempty"`
}

type WorkoutBlock struct {
	Exercises      []WorkoutExercise `json:"exercises"`
	RestAfterBlock int               `json:"restAfterBlock"`
}

type Workout struct {
	ID             id.UUID        `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	ExerciseBlocks []WorkoutBlock `json:"exerciseBlocks"`
	CreatedAt      time.Time      `json:"createdAt"`
	ModifiedAt     time.Time      `json:"modifiedAt"`
}

type WorkoutRequestParams struct {
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Exercises   []WorkoutBlock `json:"exerciseBlocks"`
}

func NewWorkoutService(repo repositories.WorkOutRepository, secrect string) *WorkoutService {
	return &WorkoutService{workoutRepo: repo, Secret: secrect}
}

func (ws *WorkoutService) createWorkoutParams(r WorkoutRequestParams) database.CreateWorkOutParams {

	return database.CreateWorkOutParams{
		Name:        r.Name,
		Description: NoneNullToNullString(r.Description),
	}
}

func (ws *WorkoutService) CreateWorkout(ctx context.Context, wrp WorkoutRequestParams) (Workout, error) {

	dbCreateParams := ws.createWorkoutParams(wrp)
	dbWO, err := ws.workoutRepo.CreateDBWO(ctx, &dbCreateParams)
	if err != nil {
		return Workout{}, err
	}

	convertedWO, err := ws.ConvertDBWorkoutToWorkout(dbWO)
	if err != nil {
		return Workout{}, err
	}

	return convertedWO, nil
}

func (ws *WorkoutService) ConvertDBWorkoutToWorkout(d database.Workout) (Workout, error) {

	// exercises := make([]WorkoutBlock, 0)
	// if d.Exercises.Valid {
	// 	err := json.Unmarshal(d.Exercises.RawMessage, &exercises)
	// 	if err != nil {
	// 		return Workout{}, err
	// 	}
	// }

	return Workout{
		ID:          d.ID,
		Name:        d.Name,
		Description: *NullStringToString(d.Description),
		//ExerciseBlocks: exercises,
		CreatedAt:  *NullTimeToTime(d.CreatedAt),
		ModifiedAt: *NullTimeToTime(d.ModifiedAt),
	}, nil
}
