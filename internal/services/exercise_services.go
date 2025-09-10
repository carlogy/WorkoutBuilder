package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
	id "github.com/google/uuid"
)

type ExerciseService struct {
	exerciseRepo repositories.ExerciseRepository
}

type ExerciseType string

const (
	ExerciseWeightedReps ExerciseType = "WeightedReps"
	ExerciseBodyWeight   ExerciseType = "BodyWeight"
	ExerciseStaticHolds  ExerciseType = "StaticHolds"
	ExerciseCardio       ExerciseType = "Cardio"
)

type Exercise struct {
	ID                    id.UUID        `json:"id,omitempty"`
	Name                  string         `json:"name"`
	ExerciseType          ExerciseType   `json:"exerciseType"`
	Equipment             string         `json:"equipment"`
	PrimaryMuscleGroups   []MuscleGroups `json:"primaryMuscleGroups"`
	SecondaryMuscleGroups []MuscleGroups `json:"secondaryMuscleGoups"`
	Description           *string        `json:"description"`
	CreatedAt             *time.Time     `json:"createdAt,omitempty"`
	ModifiedAt            *time.Time     `json:"modifiedAt,omitempty"`
}

type MuscleGroups struct {
	ID          id.UUID    `json:"id,omitempty"`
	BodyPart    string     `json:"bodyPart"`
	MuscleGroup string     `json:"muscleGroup"`
	MuscleName  string     `json:"muscleName"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	ModifiedAt  *time.Time `json:"modifiedAt,omitempty"`
}

type ExerciseMuscleGroups struct {
	ID                   int
	ExerciseID           id.UUID
	MuscleGroupID        id.UUID
	PrimaryMuscleGroup   bool
	SecondaryMuscleGroup bool
	CreatedAt            *time.Time
	ModifiedAt           *time.Time
}

type ExerciseRequestParams struct {
	Name                  string         `json:"name"`
	ExerciseType          ExerciseType   `json:"exerciseType"`
	Equipment             string         `json:"equipment"`
	PrimaryMuscleGroups   []MuscleGroups `json:"primaryMuscleGroups"`
	SecondaryMuscleGroups []MuscleGroups `json:"secondaryMuscleGroups"`
	Description           *string        `json:"description"`
}

func NewExerciseService(repo repositories.ExerciseRepository, secret string) *ExerciseService {
	return &ExerciseService{
		exerciseRepo: repo,
	}
}

func (es *ExerciseService) createExerciseParams(erp ExerciseRequestParams) db.CreateExerciseParams {

	hasPrimary := len(erp.PrimaryMuscleGroups) > 0
	hasSecondary := len(erp.SecondaryMuscleGroups) > 0

	return db.CreateExerciseParams{
		ID:                  id.New(),
		Name:                erp.Name,
		ExerciseType:        es.getExerciseType(erp.ExerciseType),
		Equipment:           erp.Equipment,
		Description:         NoneNullToNullString(erp.Description),
		HasPrimaryMuscles:   NoneNullBoolToNullBull(hasPrimary),
		HasSecondaryMuscles: NoneNullBoolToNullBull(hasSecondary),
	}
}

func (es *ExerciseService) createExercieMuscleGroupsParams(mgID, exID id.UUID, isPrimary, isSecondary bool) db.CreateExerciseMuscleGroupsParams {
	return db.CreateExerciseMuscleGroupsParams{
		ExerciseID:      exID,
		MuscleGroupsID:  mgID,
		PrimaryMuscle:   NoneNullBoolToNullBull(isPrimary),
		SecondaryMuscle: NoneNullBoolToNullBull(isSecondary),
	}
}

func (es *ExerciseService) createMuscleGroupParams(mgp MuscleGroups) db.CreateMuscleGroupParams {
	return db.CreateMuscleGroupParams{
		ID:          id.New(),
		BodyPart:    mgp.BodyPart,
		MuscleName:  mgp.MuscleName,
		MuscleGroup: mgp.MuscleGroup,
	}
}

func (es *ExerciseService) CreateExercise(ctx context.Context, erp ExerciseRequestParams) (Exercise, error) {

	exists, err := es.exerciseRepo.CheckExerciseExists(ctx, erp.Name)
	if err != nil {
		return Exercise{}, err
	}
	if exists {
		fmt.Println("Exercise already exists: ", erp.Name)
		return Exercise{}, errors.New("exercise already exists")
	}

	dbEX := es.createExerciseParams(erp)

	mgSlice := []db.CreateMuscleGroupParams{}
	emgSlice := []db.CreateExerciseMuscleGroupsParams{}

	for _, primaryMuscle := range erp.PrimaryMuscleGroups {

		dbMG, _ := es.exerciseRepo.GetMuscleGroupByName(ctx, primaryMuscle.MuscleName)

		if dbMG.ID == id.Nil {
			mgp := es.createMuscleGroupParams(primaryMuscle)
			mgSlice = append(mgSlice, mgp)

			exMG := es.createExercieMuscleGroupsParams(mgp.ID, dbEX.ID, true, false)
			emgSlice = append(emgSlice, exMG)
			continue
		}

		exMG := es.createExercieMuscleGroupsParams(dbMG.ID, dbEX.ID, true, false)
		emgSlice = append(emgSlice, exMG)
	}

	for _, secondaryMuscle := range erp.SecondaryMuscleGroups {

		dbMG, _ := es.exerciseRepo.GetMuscleGroupByName(ctx, secondaryMuscle.MuscleName)

		if dbMG.ID == id.Nil {
			mgp := es.createMuscleGroupParams(secondaryMuscle)
			mgSlice = append(mgSlice, mgp)

			exMG := es.createExercieMuscleGroupsParams(mgp.ID, dbEX.ID, false, true)
			emgSlice = append(emgSlice, exMG)
			continue
		}

		exMG := es.createExercieMuscleGroupsParams(dbMG.ID, dbEX.ID, false, true)
		emgSlice = append(emgSlice, exMG)
	}

	createdExercise, err := es.exerciseRepo.CreateDBExercise(ctx, dbEX, mgSlice, emgSlice)
	if err != nil {
		fmt.Println("Error with repo creating exercise: ", err)
		return Exercise{}, err
	}

	//To do: query all items and return constructed created entities

	fullExercise := Exercise{
		ID:           createdExercise.ID,
		Name:         createdExercise.Name,
		ExerciseType: ExerciseType(createdExercise.ExerciseType),
		Equipment:    createdExercise.Equipment,
		Description:  NullStringToString(createdExercise.Description),
		CreatedAt:    NullTimeToTime(createdExercise.CreatedAt),
		ModifiedAt:   NullTimeToTime(createdExercise.ModifiedAt),
	}

	exMuscleGroups, err := es.exerciseRepo.GetMuscleGroupsByExerciseID(ctx, createdExercise.ID)
	if err != nil {
		fmt.Println("Error getting muscle groups by exerciseID: ", err)
	}

	for _, muscleGroup := range exMuscleGroups {

		if NullBoolToBool(muscleGroup.ExerciseMuscleGroup.PrimaryMuscle) {
			fullExercise.PrimaryMuscleGroups = append(fullExercise.PrimaryMuscleGroups, MuscleGroups{ID: muscleGroup.MuscleGroup.ID, BodyPart: muscleGroup.MuscleGroup.BodyPart, MuscleGroup: muscleGroup.MuscleGroup.MuscleGroup, MuscleName: muscleGroup.MuscleGroup.MuscleName, CreatedAt: NullTimeToTime(muscleGroup.MuscleGroup.CreatedAt),
				ModifiedAt: NullTimeToTime(muscleGroup.MuscleGroup.ModifiedAt),
			})
			continue
		}

		if NullBoolToBool(muscleGroup.ExerciseMuscleGroup.SecondaryMuscle) {
			fullExercise.SecondaryMuscleGroups = append(fullExercise.SecondaryMuscleGroups,
				MuscleGroups{ID: muscleGroup.MuscleGroup.ID, BodyPart: muscleGroup.MuscleGroup.BodyPart, MuscleGroup: muscleGroup.MuscleGroup.MuscleGroup, MuscleName: muscleGroup.MuscleGroup.MuscleName, CreatedAt: NullTimeToTime(muscleGroup.MuscleGroup.CreatedAt),
					ModifiedAt: NullTimeToTime(muscleGroup.MuscleGroup.ModifiedAt)})
		}
	}

	fmt.Println(fullExercise)

	return fullExercise, nil
}

func (es *ExerciseService) GetFullExerciseByID(ctx context.Context, exId string) {

	// TO do query and construct Full exercise
}

func (es *ExerciseService) isValidExerciseType(t ExerciseType) bool {

	switch t {
	case ExerciseWeightedReps, ExerciseBodyWeight, ExerciseStaticHolds, ExerciseCardio:
		{
			return true
		}
	default:
		return false
	}
}

func (es *ExerciseService) getExerciseType(t ExerciseType) string {

	if es.isValidExerciseType(t) {
		return string(t)
	}
	return "WeightedReps"
}

func (es *ExerciseService) ConvertDBexerciseToExercise(e db.Exercise) Exercise {
	exercise := Exercise{
		ID:           e.ID,
		Name:         e.Name,
		ExerciseType: ExerciseType(es.getExerciseType(ExerciseType(e.ExerciseType))),
		Equipment:    e.Equipment,
		// PrimaryMuscleGroups:   ConvertRawJSONTOMap[string](e.HasPrimaryMuscles),
		// SecondaryMuscleGroups: ConvertRawJSONTOMap[string](e.SecondaryMuscleGroups),
		Description: NullStringToString(e.Description),
		CreatedAt:   NullTimeToTime(e.CreatedAt),
		ModifiedAt:  NullTimeToTime(e.ModifiedAt),
	}

	return exercise
}

// func ConvertStringToUUID(s string) (id.UUID, error) {
// 	err := id.Validate(s)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}

// 	newUUID, err := id.Parse(s)
// 	if err != nil {
// 		return id.Nil, err
// 	}

// 	return newUUID, nil
// }
