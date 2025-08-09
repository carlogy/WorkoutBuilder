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

func createUpdateUserExerciseParams(jsonParams CreateUserExerciseParams, recordID uuid.UUID) (database.UpdateUserExcerciseRecordByIdParams, error) {

	jsonSetWeight, err := services.ConvertMapsToRawJSON(jsonParams.SetsWeight)
	if err != nil {
		fmt.Println("Error converting sets weight to rawjson ", err)
		return database.UpdateUserExcerciseRecordByIdParams{}, err
	}

	dbUE := database.UpdateUserExcerciseRecordByIdParams{
		SetsWeight:     jsonSetWeight,
		Rest:           services.NoneNullIntToNullInt(jsonParams.Rest),
		Duration:       services.NoneNullIntToNullInt(jsonParams.Duration),
		DeclineIncline: services.NoneNullIntToNullInt(jsonParams.Decline_Incline),
		Notes:          services.NoneNullToNullString(jsonParams.Notes),
		ID:             recordID,
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

func (ueh *UserExerciseHandler) GetUserExerciseHandler(w http.ResponseWriter, r *http.Request) {

	recordID := r.PathValue("id")
	recordUUID, err := uuid.Parse(recordID)
	if err != nil {
		http.Error(w, "Unauthorized uuid format", http.StatusNotFound)
		return
	}

	tokenUserid, err := ueh.conf.GetUserIDFromToken(r.Header)
	if err != nil {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	dbRecord, err := ueh.conf.db.GetUserExerciseRecordById(r.Context(), recordUUID)
	if err != nil {
		http.Error(w, "Record Not found", http.StatusNotFound)
		fmt.Println("Error getting record from db ", err)
		return
	}

	if dbRecord.Userid != tokenUserid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error exerciseRecordUser and request token user don't match ", err)
		return
	}

	userExercise, err := services.ConvertDBUserExerciseToUserExercise(dbRecord)
	if err != nil {
		http.Error(w, "Error preparing json reponse ", http.StatusInternalServerError)
		fmt.Println("Errorconverting db record for json response ", err)
		return
	}

	ueh.writeJSONResponse(w, userExercise, http.StatusOK)
}

func (ueh *UserExerciseHandler) UpdateUserExerciseHandler(w http.ResponseWriter, r *http.Request) {

	recordID := r.PathValue("id")
	uuidRecordID, err := uuid.Parse(recordID)
	if err != nil {
		http.Error(w, "Unauthorized uuid format", http.StatusNotFound)
		return
	}

	tokenUserid, err := ueh.conf.GetUserIDFromToken(r.Header)
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

	if tokenUserid != body.UserID {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	dbupdateParams, err := createUpdateUserExerciseParams(body, uuidRecordID)
	if err != nil {
		http.Error(w, "Error creating entity", 500)
		fmt.Println("Error creating userExercise params: ", err)
		return
	}

	updatedUERecord, err := ueh.conf.db.UpdateUserExcerciseRecordById(r.Context(), dbupdateParams)
	if err != nil {
		http.Error(w, "Error updating record", http.StatusInternalServerError)
		fmt.Println("Error updating record ", err)
		return
	}

	convertedUERecord, err := services.ConvertDBUserExerciseToUserExercise(updatedUERecord)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		fmt.Println("Error converting db to response record: ", err)
		return
	}

	ueh.writeJSONResponse(w, convertedUERecord, http.StatusOK)
}

func (ueh *UserExerciseHandler) DeleteUserExerciseRecordById(w http.ResponseWriter, r *http.Request) {
	recordID := r.PathValue("id")
	uuidRecordID, err := uuid.Parse(recordID)
	if err != nil {
		http.Error(w, "Unauthorized uuid format", http.StatusNotFound)
		return
	}

	tokenUserid, err := ueh.conf.GetUserIDFromToken(r.Header)
	if err != nil {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	dbRecord, err := ueh.conf.db.GetUserExerciseRecordById(r.Context(), uuidRecordID)
	if err != nil {
		http.Error(w, "Error invalid recordID", http.StatusInternalServerError)
		fmt.Println("Error querying for record in db: ", err)
		return
	}

	if tokenUserid != dbRecord.Userid {
		http.Error(w, "Error validation bearer token", http.StatusUnauthorized)
		return
	}

	deletedRecord, err := ueh.conf.db.DeleteUserExerciseRecordById(r.Context(), uuidRecordID)
	if err != nil {
		http.Error(w, "Error deleting record", http.StatusInternalServerError)
		fmt.Println("Error attempting to delete record: ", err)
		return
	}

	jsonReponseRecord, err := services.ConvertDBUserExerciseToUserExercise(deletedRecord)
	if err != nil {
		http.Error(w, "Error preparing jsonResponse", http.StatusInternalServerError)
		fmt.Println("Error converting db record to jsonResponse: ", err)
		return
	}

	ueh.writeJSONResponse(w, jsonReponseRecord, http.StatusOK)
}
