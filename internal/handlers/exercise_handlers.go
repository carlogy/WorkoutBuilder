package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
	id "github.com/google/uuid"
)

type ExerciseHandler struct {
	exerciseService *services.ExerciseService
	authService     *services.AuthService
}

func NewExerciseHandler(es services.ExerciseService, as services.AuthService) ExerciseHandler {
	return ExerciseHandler{exerciseService: &es, authService: &as}
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, eh.authService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	jep := services.ExerciseRequestParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jep)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		fmt.Println("Error decoing request body: ", err)
		return
	}

	ex, err := eh.exerciseService.CreateExercise(r.Context(), jep)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"unique_exercise_name\"" {
			http.Error(w, "Exercise already exists. Please use existing exercise", http.StatusUnprocessableEntity)
			fmt.Println("Exercise already created: ", err)
			return
		}
		if err.Error() == "exercise already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(500)
		fmt.Println(err.Error())
		return
	}

	eh.writeJSONResponse(w, ex, http.StatusOK)
}

func (eh *ExerciseHandler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	pathID := r.PathValue("id")
	exID, err := id.Parse(pathID)
	if err != nil {
		http.Error(w, "Erroronous workoutID received", http.StatusNotFound)
		fmt.Println("Error parsing string id to uuid: ", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, eh.authService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	e, err := eh.exerciseService.GetFullExerciseByID(r.Context(), exID)
	if err != nil {
		http.Error(w, "Error: Exercise does not exist", 500)
		fmt.Printf("Experienced error when querying db for exercise: \t%v\n", err)
		return
	}

	eh.writeJSONResponse(w, e, 200)
}

func (eh *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, eh.authService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	exList, err := eh.exerciseService.GetAllExercises(r.Context())
	if err != nil {
		fmt.Println("Error getting list of exercises from service: ", err)
		http.Error(w, "Error getting list of exercises", http.StatusInternalServerError)
		return
	}

	eh.writeJSONResponse(w, exList, 200)

}

func (eh *ExerciseHandler) DeleteExerciseByID(w http.ResponseWriter, r *http.Request) {

	pathID := r.PathValue("id")
	exID, err := id.Parse(pathID)
	if err != nil {
		http.Error(w, "Erroronous workoutID received", http.StatusNotFound)
		fmt.Println("Error parsing string id to uuid: ", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, eh.authService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	deletedEx, err := eh.exerciseService.DeleteExerciseByID(r.Context(), exID)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		fmt.Printf("Experienced err while service deleted exercise:\t%v\n", err)
		return
	}

	eh.writeJSONResponse(w, deletedEx, 200)
}
