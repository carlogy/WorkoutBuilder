package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
	"github.com/google/uuid"
)

type UserExerciseHandler struct {
	conf *ApiConfig
}

type CreateUserExerciseParams struct {
	UserID          uuid.UUID              `json:"userID"`
	ExerciseID      uuid.UUID              `json:"exerciseID"`
	SetsWeight      map[string]map[int]int `json:"sets_weight"`
	Rest            *int                   `json:"rest"`
	Duration        *int                   `json:"durantion"`
	Decline_Incline *int                   `json:"decline_incline"`
	Notes           *string                `json:"notes"`
}

func NewUserExerciseHanlder(ac *ApiConfig) UserExerciseHandler {
	return UserExerciseHandler{conf: ac}
}

func (ueh *UserExerciseHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Conten-Type", "application/json")

	dat, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Error encoding response", 500)
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

func createDBUserExerciseParams(jsonUserParam CreateUserExerciseParams) (database.CreateUserExerciseParams, error) {

	jsonSetsWeight, err := services.ConvertMapsToRawJSON(jsonUserParam.SetsWeight)
	if err != nil {
		fmt.Println("Error converting sets weight to rawjson ", err)
		return database.CreateUserExerciseParams{}, err
	}
	dbUE := database.CreateUserExerciseParams{
		Userid:         jsonUserParam.UserID,
		Exerciseid:     jsonUserParam.ExerciseID,
		SetsWeight:     jsonSetsWeight,
		Rest:           services.NoneNullIntToNullInt(jsonUserParam.Rest),
		Duration:       services.NoneNullIntToNullInt(jsonUserParam.Duration),
		DeclineIncline: services.NoneNullIntToNullInt(jsonUserParam.Decline_Incline),
		Notes:          services.NoneNullToNullString(jsonUserParam.Notes),
	}

	return dbUE, nil

}

func (ueh *UserExerciseHandler) CreateUserExerciseHandler(w http.ResponseWriter, r *http.Request) {

	id, err := ueh.conf.GetUserIDFromToken(r.Header)
	if err != nil {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)

	body := CreateUserExerciseParams{}
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	if id != body.UserID {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	dbUE, err := createDBUserExerciseParams(body)
	if err != nil {
		http.Error(w, "Error creating entity", 500)
		fmt.Println("Error creating userExercise params: ", err)
		return
	}

	createdUE, err := ueh.conf.db.CreateUserExercise(r.Context(), dbUE)
	if err != nil {
		http.Error(w, "Error logging exercise", http.StatusInternalServerError)
		fmt.Println("Error writing UserExercise to db ", err)
		return
	}

	userExercise, err := services.ConvertDBUserExerciseToUserExercise(createdUE)
	if err != nil {
		http.Error(w, "Error encoding recorded exercise", http.StatusInternalServerError)
	}

	ueh.writeJSONResponse(w, userExercise, http.StatusCreated)
}
