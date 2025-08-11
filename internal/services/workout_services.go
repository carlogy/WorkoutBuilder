package services

import (
	"encoding/json"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type Set struct {
	Weight int `json:"weight"`
	Reps   int `json:"reps"`
}

type WorkoutExercise struct {
	ExerciseID id.UUID `json:"exerciseID"`
	Sets       []Set   `json:"sets"`
}

type Workout struct {
	ID          id.UUID           `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises"`
	CreatedAt   time.Time         `json:"createdAt"`
	ModifiedAt  time.Time         `json:"modifiedAt"`
}

func ConvertDBWorkoutToWorkout(d database.Workout) (Workout, error) {

	exercises := make([]WorkoutExercise, 0)
	if d.Exercises.Valid {
		err := json.Unmarshal(d.Exercises.RawMessage, &exercises)
		if err != nil {
			return Workout{}, err
		}
	}

	return Workout{
		ID:          d.ID,
		Name:        d.Name,
		Description: *NullStringToString(d.Description),
		Exercises:   exercises,
		CreatedAt:   *NullTimeToTime(d.CreatedAt),
		ModifiedAt:  *NullTimeToTime(d.ModifiedAt),
	}, nil

}
