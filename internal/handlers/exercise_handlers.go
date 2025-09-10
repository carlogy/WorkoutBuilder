package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
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

	// fmt.Println(jep)

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

func (eh *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {

	//To Do: refactor using handler, service, repository implementation and normalizaion of tables

	// exerciseList, err := eh.conf.db.GetExercises(r.Context())
	// if err != nil {
	// 	http.Error(w, "Error: getting exercises", 500)
	// 	fmt.Printf("Error querying for exercises: %v", err)
	// 	return

	// }

	// type jsonBodyResp struct {
	// 	Exercises []services.Exercise `json:"exercises"`
	// 	Total     int                 `json:"total"`
	// }

	// updatedExerciseList := make([]services.Exercise, 0)
	// for _, e := range exerciseList {
	// 	exercise := services.ConvertDBexerciseToExercise(e)
	// 	updatedExerciseList = append(updatedExerciseList, exercise)
	// }
	// jbr := jsonBodyResp{Exercises: updatedExerciseList, Total: len(updatedExerciseList)}

	// eh.writeJSONResponse(w, jbr, 200)

}

func (eh *ExerciseHandler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	//To Do: refactor using handler, service, repository implementation and normalizaion of tables

	// id := r.PathValue("id")
	// if id == "" {

	// 	http.Error(w, "Error: getting path value", 500)
	// 	fmt.Printf("Id string received:\t%v", id)
	// 	return
	// }

	// uuid, err := services.ConvertStringToUUID(id)
	// if err != nil {
	// 	http.Error(w, "Invalid id passed", 500)
	// 	fmt.Println("Error converting pathVal to UUID", err)
	// 	return
	// }

	// dbEx, err := eh.conf.db.GetExerciseById(r.Context(), uuid)
	// if err != nil {
	// 	http.Error(w, "Error: Exercise does not exist", 500)
	// 	fmt.Printf("Experienced error when querying db for exercise: \t%v\n", err)
	// 	return
	// }

	// exercise := services.ConvertDBexerciseToExercise(dbEx)

	// eh.writeJSONResponse(w, exercise, 200)
}

func (eh *ExerciseHandler) DeleteExerciseByID(w http.ResponseWriter, r *http.Request) {

	//To Do: refactor using handler, service, repository implementation and normalizaion of tables

	// id := r.PathValue("id")
	// if id == "" {
	// 	http.Error(w, "Must provide exercise Id to delete", 500)
	// 	fmt.Println("Invalid id passed:\t", id)
	// 	return
	// }

	// uuid, err := services.ConvertStringToUUID(id)
	// if err != nil {
	// 	http.Error(w, "Invalid id passed", 500)
	// 	fmt.Println("Error converting pathVal to UUID", err)
	// 	return
	// }

	// deletedEx, err := eh.conf.db.DeleteExerciseById(r.Context(), uuid)
	// if err != nil {
	// 	http.Error(w, "Exercise not found", 404)
	// 	fmt.Printf("Experienced err while deleting exercise:\t%v\n", err)
	// 	return
	// }

	// ex := services.ConvertDBexerciseToExercise(deletedEx)

	// eh.writeJSONResponse(w, ex, 200)
}
