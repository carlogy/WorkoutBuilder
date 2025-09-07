package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/services"
)

type WorkoutHandler struct {
	workoutService *services.WorkoutService
	authService    *services.AuthService
}

type WorkoutParamsRequest struct {
	Name        string                  `json:"name"`
	Description *string                 `json:"description"`
	Exercises   []services.WorkoutBlock `json:"exerciseBlocks"`
}

func NewWorkoutHandler(ws services.WorkoutService, as services.AuthService) WorkoutHandler {
	return WorkoutHandler{workoutService: &ws, authService: &as}
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

// to do get rid of this function
// func createWorkoutDBParams(r WorkoutParamsRequest) (database.CreateWorkOutParams, error) {

// 	// we := make([]services.WorkoutExercise, 0)
// 	// for i, e := range r.Exercises {
// 	// 	we = append(we, e.Exercises[i])
// 	// }

// 	// exercises := pqtype.NullRawMessage{}
// 	// dat, err := json.Marshal(we)
// 	// if err != nil {
// 	// 	return database.CreateWorkOutParams{}, err
// 	// }
// 	// exercises.RawMessage = dat
// 	// exercises.Valid = true

// 	return database.CreateWorkOutParams{
// 		Name:        r.Name,
// 		Description: services.NoneNullToNullString(r.Description),
// 	}, nil
// }

func (wh WorkoutHandler) CreateWorkoutHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, wh.workoutService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	jwp := services.WorkoutRequestParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jwp)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		fmt.Println("Error decoing request body: ", err)
		return
	}

	workout, err := wh.workoutService.CreateWorkout(r.Context(), jwp)
	if err != nil {
		http.Error(w, "Error writing to db", http.StatusInternalServerError)
		fmt.Println("Error writing to db: ", err)
		return
	}

	wh.writeJSONResponse(w, workout, http.StatusOK)

}

func (wh *WorkoutHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {

	//to do re-implement using service and repositories

	// token, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
	// 	fmt.Println("Error gettting auth token: ", err)
	// 	return
	// }

	// _, err = auth.ValidateJWT(token, wh.workoutService.Secret)
	// if err != nil {
	// 	http.Error(w, "Invalid Token", http.StatusUnauthorized)
	// 	fmt.Println("Error validating JWT token: ", err)
	// 	return
	// }

	// dbWorkouts, err := wh.conf.db.GetWorkouts(r.Context())
	// if err != nil {
	// 	http.Error(w, "Error getting workouts", http.StatusInternalServerError)
	// 	fmt.Println("Error querying all workouts: ", err)
	// 	return
	// }

	// type jsonbody struct {
	// 	Workouts []services.Workout `json:"workouts"`
	// 	Total    int                `json:"total"`
	// }

	// convertedWorkoutList := make([]services.Workout, 0)
	// for _, v := range dbWorkouts {
	// 	cw, err := services.ConvertDBWorkoutToWorkout(v)
	// 	if err != nil {
	// 		fmt.Println("Unable to covert workout: ", err)
	// 		continue
	// 	}
	// 	convertedWorkoutList = append(convertedWorkoutList, cw)
	// }

	// jbody := jsonbody{Workouts: convertedWorkoutList, Total: len(convertedWorkoutList)}

	// w.Header().Set("Conten-Type", "application/json")

	// dat, err := json.Marshal(jbody)

	// if err != nil {
	// 	http.Error(w, "Error serializing response", 500)
	// 	fmt.Printf("Error marshalling db response: %v", err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
	// _, err = w.Write(dat)
	// if err != nil {
	// 	log.Printf("Error writing response: %v\n", err)
	// 	return
	// }
}

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {

	//to do re-implement using service and repositories

	// pathID := r.PathValue("id")
	// wuuid, err := uuid.Parse(pathID)
	// if err != nil {
	// 	http.Error(w, "Erroronous workoutID received", http.StatusInternalServerError)
	// 	fmt.Println("Error parsing string id to uuid: ", err)
	// 	return
	// }

	// token, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
	// 	fmt.Println("Error gettting auth token: ", err)
	// 	return
	// }

	// _, err = auth.ValidateJWT(token, wh.conf.secret)
	// if err != nil {
	// 	http.Error(w, "Invalid Token", http.StatusUnauthorized)
	// 	fmt.Println("Error validating JWT token: ", err)
	// 	return
	// }

	// dbWorkout, err := wh.conf.db.GetWorkoutByID(r.Context(), wuuid)
	// if err != nil {
	// 	http.Error(w, "Workout not found", http.StatusNotFound)
	// 	fmt.Println("Error querying db for workout: ", err)
	// 	return
	// }

	// convertedWorkout, err := services.ConvertDBWorkoutToWorkout(dbWorkout)
	// if err != nil {
	// 	http.Error(w, "Error preparing json response", http.StatusInternalServerError)
	// 	fmt.Println("Error marshalling db workout to json response ")
	// 	return
	// }

	// wh.writeJSONResponse(w, convertedWorkout, http.StatusOK)
}

func (wh *WorkoutHandler) DeleteWorkoutById(w http.ResponseWriter, r *http.Request) {

	//to do re-implement using service and repositories

	// pathID := r.PathValue("id")
	// wuuid, err := uuid.Parse(pathID)
	// if err != nil {
	// 	http.Error(w, "Erroronous workoutID received", http.StatusInternalServerError)
	// 	fmt.Println("Error parsing string id to uuid: ", err)
	// 	return
	// }

	// token, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
	// 	fmt.Println("Error gettting auth token: ", err)
	// 	return
	// }

	// _, err = auth.ValidateJWT(token, wh.conf.secret)
	// if err != nil {
	// 	http.Error(w, "Invalid Token", http.StatusUnauthorized)
	// 	fmt.Println("Error validating JWT token: ", err)
	// 	return
	// }

	// dbWorkout, err := wh.conf.db.DeleteWorkoutByID(r.Context(), wuuid)
	// if err != nil {
	// 	http.Error(w, "Workout not found", http.StatusNotFound)
	// 	fmt.Println("Error querying db for workout: ", err)
	// 	return
	// }

	// convertedWorkout, err := services.ConvertDBWorkoutToWorkout(dbWorkout)
	// if err != nil {
	// 	http.Error(w, "Error preparing json response", http.StatusInternalServerError)
	// 	fmt.Println("Error marshalling db workout to json response ")
	// 	return
	// }

	// wh.writeJSONResponse(w, convertedWorkout, http.StatusOK)
}
