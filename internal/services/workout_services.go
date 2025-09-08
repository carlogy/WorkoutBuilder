package services

import (
	"context"
	"fmt"
	"time"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
	id "github.com/google/uuid"
)

type WorkoutService struct {
	workoutRepo repositories.WorkOutRepository
	Secret      string
}

type Set struct {
	ID               id.UUID `json:"id,omitempty"`
	Ordinal          int     `json:"ordinal"`
	Weight           float64 `json:"weight"`
	Reps             int     `json:"reps"`
	StaticHolds      int     `json:"statichHoldTime"`
	WokoutExerciseID id.UUID `json:"workoutExerciseID"`
}

type WorkoutExercise struct {
	ID         id.UUID `json:"workoutExerciseID,omitempty"`
	Ordinal    int     `json:"ordinal"`
	ExerciseID id.UUID `json:"exerciseID"`
	Sets       []Set   `json:"sets"`
	Notes      string  `JSON:"notes,omitempty"`
}

type WorkoutBlock struct {
	ID             id.UUID           `json:"workoutBlockID,omitempty"`
	Ordinal        int               `json:"ordinal"`
	Exercises      []WorkoutExercise `json:"exercises"`
	RestAfterBlock int               `json:"restAfterBlock"`
	WorkoutID      id.UUID           `json:"workoutID"`
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

func (ws *WorkoutService) createWorkoutParams(r WorkoutRequestParams) db.CreateWorkOutParams {

	return db.CreateWorkOutParams{
		ID:          id.New(),
		Name:        r.Name,
		Description: NoneNullToNullString(r.Description),
	}
}

func (ws *WorkoutService) createWorkoutBlockParams(wob WorkoutBlock, woID id.UUID) db.CreateWorkoutBlocksParams {
	return db.CreateWorkoutBlocksParams{
		ID:                    wob.ID,
		Ordinal:               NoneNullIntToNullInt32(&wob.Ordinal),
		Workoutid:             woID,
		RestsecondsAfterBlock: NoneNullIntToNullInt32(&wob.RestAfterBlock),
	}
}

func (ws *WorkoutService) createExerciseSets(workoutExerciseID id.UUID, set Set) db.CreateExerciseSetsParams {

	ordinalStore := NoneNullIntToNullInt32(&set.Ordinal)

	return db.CreateExerciseSetsParams{
		ID:                id.New(),
		Ordinal:           ordinalStore.Int32,
		WorkoutExerciseid: workoutExerciseID,
		Weight:            set.Weight,
		Reps:              NoneNullIntToNullInt32(&set.Reps),
		StaticHoldTime:    NoneNullIntToNullInt32(&set.StaticHolds),
	}
}

func (ws *WorkoutService) createWorkoutExercisesParams(woe WorkoutExercise, woBlockID id.UUID) db.CreateWorkoutExercisesParams {
	return db.CreateWorkoutExercisesParams{
		ID:             woe.ID,
		Ordinal:        NoneNullIntToNullInt32(&woe.Ordinal),
		WorkoutBlockid: woBlockID,
		Exerciseid:     woe.ExerciseID,
		Notes:          NoneNullToNullString(&woe.Notes),
	}
}

func (ws *WorkoutService) CreateWorkout(ctx context.Context, wrp WorkoutRequestParams) (Workout, error) {

	dbCreateParams := ws.createWorkoutParams(wrp)

	woBlockSlice := make([]db.CreateWorkoutBlocksParams, 0)
	woExSlice := make([]db.CreateWorkoutExercisesParams, 0)
	woExSets := make([]db.CreateExerciseSetsParams, 0)

	for _, woBlock := range wrp.Exercises {

		woBlock.ID = id.New()
		woBlockSlice = append(woBlockSlice, ws.createWorkoutBlockParams(woBlock, dbCreateParams.ID))

		for _, woEx := range woBlock.Exercises {

			woEx.ID = id.New()
			woExSlice = append(woExSlice, ws.createWorkoutExercisesParams(woEx, woBlock.ID))

			for _, set := range woEx.Sets {

				woExSets = append(woExSets, ws.createExerciseSets(woEx.ID, set))
			}
		}
	}

	dbWO, err := ws.workoutRepo.CreateDBWO(ctx, &dbCreateParams, woBlockSlice, woExSlice, woExSets)
	if err != nil {
		return Workout{}, err
	}

	// think about using helper functions, instead using getWorkoutID service method ...?

	convertedWO, err := ws.GetWorkoutByID(ctx, dbWO.ID)
	if err != nil {
		return Workout{}, err
	}

	return convertedWO, nil
}

func (ws *WorkoutService) GetWorkoutByID(ctx context.Context, woID id.UUID) (Workout, error) {

	w, err := ws.workoutRepo.GetDBWOByID(ctx, woID)
	if err != nil {
		fmt.Println("Error getting workout by ID from repo: ", err)
		return Workout{}, err
	}

	convertedWO, err := ws.ConvertDBWorkoutToWorkout(w)
	if err != nil {
		return Workout{}, err
	}

	// to do break out into helper functions
	woExerciseBlockMap := map[*WorkoutExercise]id.UUID{}
	woSetExerciseMap := map[Set]id.UUID{}

	dbWoBlocks, err := ws.workoutRepo.GetDBWOBlocksByWOID(ctx, convertedWO.ID)
	if err != nil {
		fmt.Println("Error getting back created workouts: ", err)
		return Workout{}, err
	}

	dbWoExercises, err := ws.workoutRepo.GetDBWOExByWOID(ctx, convertedWO.ID)
	if err != nil {
		fmt.Println("Error getting back created exercises: ", err)
		return Workout{}, err
	}

	dbExerciseSets, err := ws.workoutRepo.GetExSetsByWOID(ctx, convertedWO.ID)
	if err != nil {
		fmt.Println("Error getting created exercise sets: ", err)
		return Workout{}, err
	}

	for _, set := range dbExerciseSets {

		s := Set{
			ID:               set.ID,
			Ordinal:          int(set.Ordinal),
			Weight:           set.Weight,
			Reps:             int(set.Reps.Int32),
			StaticHolds:      int(set.StaticHoldTime.Int32),
			WokoutExerciseID: set.WorkoutExerciseid,
		}

		woSetExerciseMap[s] = s.WokoutExerciseID
	}

	for _, exercise := range dbWoExercises {

		e := WorkoutExercise{
			ID:         exercise.ID,
			Ordinal:    int(exercise.Ordinal.Int32),
			Notes:      exercise.Notes.String,
			ExerciseID: exercise.Exerciseid,
		}

		for set, woEXID := range woSetExerciseMap {

			if woEXID != e.ID {
				continue
			}

			e.Sets = append(e.Sets, set)
		}

		woExerciseBlockMap[&e] = exercise.WorkoutBlockid
	}

	for _, block := range dbWoBlocks {

		b := WorkoutBlock{
			ID:             block.ID,
			Ordinal:        int(block.Ordinal.Int32),
			RestAfterBlock: int(block.RestsecondsAfterBlock.Int32),
			WorkoutID:      block.Workoutid,
		}

		for exercise, woBlockID := range woExerciseBlockMap {

			if woBlockID != b.ID {
				continue
			}

			b.Exercises = append(b.Exercises, *exercise)

		}

		convertedWO.ExerciseBlocks = append(convertedWO.ExerciseBlocks, b)
	}

	return convertedWO, nil
}

func (ws *WorkoutService) ConvertDBWorkoutToWorkout(d db.Workout) (Workout, error) {

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
