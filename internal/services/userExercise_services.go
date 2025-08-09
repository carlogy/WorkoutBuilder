package services

import (
	"encoding/json"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/google/uuid"
)

type UserExercise struct {
	ID              uuid.UUID              `json:"recordID"`
	UserID          uuid.UUID              `json:"userID"`
	ExerciseID      uuid.UUID              `json:"exerciseID"`
	SetsWeight      map[string]map[int]int `json:"sets_weight"`
	Rest            *int                   `json:"rest"`
	Duration        *int                   `json:"durantion"`
	Decline_Incline *int                   `json:"decline_incline"`
	Notes           *string                `json:"notes"`
	CreatedAt       *time.Time             `json:"createdAt"`
	ModifiedAt      *time.Time             `json:"modifiedAt"`
}

func ConvertDBUserExerciseToUserExercise(ue database.UserExercise) (UserExercise, error) {

	var setsWeight map[string]map[int]int
	if ue.SetsWeight.Valid {
		err := json.Unmarshal(ue.SetsWeight.RawMessage, &setsWeight)
		if err != nil {
			return UserExercise{}, err
		}
	}

	userExercise := UserExercise{
		ID:         ue.ID,
		UserID:     ue.Userid,
		ExerciseID: ue.Exerciseid,
		SetsWeight: setsWeight,
		Rest:       NullInttoInt(ue.Rest),
		Duration:   NullInttoInt(ue.Duration),
		Notes:      NullStringToString(ue.Notes),
		CreatedAt:  NullTimeToTime(ue.CreatedAt),
		ModifiedAt: NullTimeToTime(ue.ModifiedAt),
	}

	return userExercise, nil
}
