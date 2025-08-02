package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
)

type ExerciseHandler struct {
	db *database.Queries
}

func NewExerciseHandler(db *database.Queries) ExerciseHandler {
	eh := ExerciseHandler{db: db}
	return eh
}

func (eh *ExerciseHandler) GetExerciseType(t services.ExerciseType) string {

	exerciseType := strconv.Itoa(int(t))
	return exerciseType
}

func (eh *ExerciseHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Conten-Type", "application/json")

	dat, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error marshalling db response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)

}

func (eh *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	type jsonExercise struct {
		Name                  string            `json:"name"`
		ExerciseType          int               `json:"exerciseType"`
		Equipment             string            `json:"equipment"`
		PrimaryMuscleGroups   map[string]string `json:"primaryMuscleGroups"`
		SecondaryMuscleGroups map[string]string `json:"secondaryMuscleGroups"`
		Description           string            `json:"description"`
		PreviousSetRepCount   map[int]int       `json:"previousSetRepCount"`
		RestBetweenSets       int               `json:"restBetweenSets"`
	}

	body := jsonExercise{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	ex, err := eh.db.CreateExercise(r.Context(), database.CreateExerciseParams{
		Name:                  body.Name,
		ExerciseType:          eh.GetExerciseType(services.ExerciseType(body.ExerciseType)),
		Equipment:             body.Equipment,
		PrimaryMuscleGroups:   services.ConvertMapToRawJSON(body.PrimaryMuscleGroups),
		SecondaryMuscleGroups: services.ConvertMapToRawJSON(body.SecondaryMuscleGroups),
		Description:           services.NoneNullToNullString(body.Description),
	},
	)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	createdExcercise := services.ConvertDBexerciseToExercise(ex)

	eh.writeJSONResponse(w, createdExcercise, 201)

}

func (eh *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {

	exerciseList, err := eh.db.GetExercises(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error getting exercises"))
		fmt.Printf("Error querying for exercises: %v", err)
		return

	}

	updatedExerciseList := make([]services.Exercise, 0)
	for i, e := range exerciseList {
		exercise := services.ConvertDBexerciseToExercise(e)
		updatedExerciseList = append(updatedExerciseList, exercise)
		fmt.Println(i, e.ID)
	}

	eh.writeJSONResponse(w, updatedExerciseList, 200)

}

func (eh *ExerciseHandler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(500)
		w.Write([]byte("Error getting  path value"))
		fmt.Printf("Id string received:\t%v", id)
		return
	}

	uuid, err := services.ConvertStringToUUID(id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Invalid id passed"))
		fmt.Println("Error converting pathVal to UUID", err)
	}

	dbEx, err := eh.db.GetExerciseById(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(500)

		w.Write([]byte("Error: Exercise does not exist"))
		fmt.Printf("Experienced error when querying db for exercise: \t%v\n", err)
		return
	}

	exercise := services.ConvertDBexerciseToExercise(dbEx)

	eh.writeJSONResponse(w, exercise, 200)
}

func (eh *ExerciseHandler) DeleteExerciseByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(500)
		w.Write([]byte("Must provide exercise Id to delete"))
	}

	uuid, err := services.ConvertStringToUUID(id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Invalid id passed"))
		fmt.Println("Error converting pathVal to UUID", err)
	}

	deletedEx, err := eh.db.DeleteExerciseById(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Exercise not found"))
		fmt.Printf("Experienced err while deleting exercise:\t%v\n", err)
	}

	ex := services.ConvertDBexerciseToExercise(deletedEx)

	eh.writeJSONResponse(w, ex, 200)
}
