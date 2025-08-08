package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/carlogy/WorkoutBuilder/internal/database"

	services "github.com/carlogy/WorkoutBuilder/internal/services"
)

type ExerciseHandler struct {
	conf *ApiConfig
}

func NewExerciseHandler(c *ApiConfig) ExerciseHandler {
	eh := ExerciseHandler{conf: c}
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
		http.Error(w, "Error serializing response", 500)
		fmt.Printf("Error marshalling db response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
		return
	}
}

func (eh *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	type jsonExercise struct {
		Name                  string            `json:"name"`
		ExerciseType          int               `json:"exerciseType"`
		Equipment             string            `json:"equipment"`
		PrimaryMuscleGroups   map[string]string `json:"primaryMuscleGroups"`
		SecondaryMuscleGroups map[string]string `json:"secondaryMuscleGroups"`
		Description           *string           `json:"description"`
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

	ex, err := eh.conf.db.CreateExercise(r.Context(), database.CreateExerciseParams{
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

	exerciseList, err := eh.conf.db.GetExercises(r.Context())
	if err != nil {
		http.Error(w, "Error: getting exercises", 500)
		fmt.Printf("Error querying for exercises: %v", err)
		return

	}

	updatedExerciseList := make([]services.Exercise, 0)
	for _, e := range exerciseList {
		exercise := services.ConvertDBexerciseToExercise(e)
		updatedExerciseList = append(updatedExerciseList, exercise)
	}

	eh.writeJSONResponse(w, updatedExerciseList, 200)

}

func (eh *ExerciseHandler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {

		http.Error(w, "Error: getting path value", 500)
		fmt.Printf("Id string received:\t%v", id)
		return
	}

	uuid, err := services.ConvertStringToUUID(id)
	if err != nil {
		http.Error(w, "Invalid id passed", 500)
		fmt.Println("Error converting pathVal to UUID", err)
		return
	}

	dbEx, err := eh.conf.db.GetExerciseById(r.Context(), uuid)
	if err != nil {
		http.Error(w, "Error: Exercise does not exist", 500)
		fmt.Printf("Experienced error when querying db for exercise: \t%v\n", err)
		return
	}

	exercise := services.ConvertDBexerciseToExercise(dbEx)

	eh.writeJSONResponse(w, exercise, 200)
}

func (eh *ExerciseHandler) DeleteExerciseByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Must provide exercise Id to delete", 500)
		fmt.Println("Invalid id passed:\t", id)
		return
	}

	uuid, err := services.ConvertStringToUUID(id)
	if err != nil {
		http.Error(w, "Invalid id passed", 500)
		fmt.Println("Error converting pathVal to UUID", err)
		return
	}

	deletedEx, err := eh.conf.db.DeleteExerciseById(r.Context(), uuid)
	if err != nil {
		http.Error(w, "Exercise not found", 404)
		fmt.Printf("Experienced err while deleting exercise:\t%v\n", err)
		return
	}

	ex := services.ConvertDBexerciseToExercise(deletedEx)

	eh.writeJSONResponse(w, ex, 200)
}
