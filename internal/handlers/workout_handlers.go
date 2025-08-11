package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/services"
	"github.com/sqlc-dev/pqtype"
)

type WorkoutHandler struct {
	conf ApiConfig
}

type jsonWorkoutParams struct {
	Name        string                     `json:"name"`
	Description *string                    `json:"description"`
	Exercises   []services.WorkoutExercise `json:"exercises"`
}

func NewWorkoutHandler(c *ApiConfig) WorkoutHandler {
	return WorkoutHandler{conf: *c}
}

func (wh *WorkoutHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
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

func createWorkoutDBParams(r jsonWorkoutParams) (database.CreateWorkOutParams, error) {

	exercises := pqtype.NullRawMessage{}
	dat, err := json.Marshal(r.Exercises)
	if err != nil {
		return database.CreateWorkOutParams{}, err
	}
	exercises.RawMessage = dat
	exercises.Valid = true

	return database.CreateWorkOutParams{
		Name:        r.Name,
		Description: services.NoneNullToNullString(r.Description),
		Exercises:   exercises,
	}, nil
}

func (wh WorkoutHandler) CreateWorkoutHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, wh.conf.secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	jwp := jsonWorkoutParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jwp)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		fmt.Println("Error decoing request body: ", err)
		return
	}

	dbCreateWorkoutParams, err := createWorkoutDBParams(jwp)
	if err != nil {
		http.Error(w, "Error saving record", http.StatusInternalServerError)
		fmt.Println("Error creating db params for insert: ", err)
		return
	}

	dbWorkout, err := wh.conf.db.CreateWorkOut(r.Context(), dbCreateWorkoutParams)
	if err != nil {
		http.Error(w, "Error writing to db", http.StatusInternalServerError)
		fmt.Println("Error writing to db: ", err)
		return
	}

	convertedWorkout, err := services.ConvertDBWorkoutToWorkout(dbWorkout)
	if err != nil {
		http.Error(w, "Error marshalling workout", http.StatusInternalServerError)
		fmt.Println("Error marshalling db workout to workout: ", err)
		return
	}

	wh.writeJSONResponse(w, convertedWorkout, http.StatusOK)

}
