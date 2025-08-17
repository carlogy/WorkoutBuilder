package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/google/uuid"
)

type ExerciseType int

const (
	WeightsWithReps = iota
	BodyWeightReps
	StaticHolds
	Cardio
)

type Exercise struct {
	ID                    uuid.UUID         `json:"id"`
	Name                  string            `json:"name"`
	ExerciseType          ExerciseType      `json:"exerciseType"`
	Equipment             string            `json:"equipment"`
	PrimaryMuscleGroups   map[string]string `json:"primaryMuscleGroups"`
	SecondaryMuscleGroups map[string]string `json:"secondaryMuscleGroups"`
	Description           *string           `json:"description"`
	CreatedAt             *time.Time        `json:"createdAt"`
	ModifiedAt            *time.Time        `json:"modifiedAt"`
}

func ConvertExerciseTypeEnum(exerciseTypeString string) ExerciseType {

	i, err := strconv.Atoi(exerciseTypeString)
	if err != nil {
		fmt.Println(err)

	}
	return ExerciseType(i)
}

func ConvertDBexerciseToExercise(e database.Exercise) Exercise {
	exercise := Exercise{
		ID:                    e.ID,
		Name:                  e.Name,
		ExerciseType:          ConvertExerciseTypeEnum(e.ExerciseType),
		Equipment:             e.Equipment,
		PrimaryMuscleGroups:   ConvertRawJSONTOMap[string](e.PrimaryMuscleGroups),
		SecondaryMuscleGroups: ConvertRawJSONTOMap[string](e.SecondaryMuscleGroups),
		Description:           NullStringToString(e.Description),
		CreatedAt:             NullTimeToTime(e.CreatedAt),
		ModifiedAt:            NullTimeToTime(e.ModifiedAt),
	}

	return exercise
}

func ConvertStringToUUID(s string) (uuid.UUID, error) {
	err := uuid.Validate(s)
	if err != nil {
		return uuid.Nil, err
	}

	newUUID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}

	return newUUID, nil
}
