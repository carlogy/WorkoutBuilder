package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/services"
	id "github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandler(us services.UserService, as services.AuthService) UserHandler {
	return UserHandler{userService: &us, authService: &as}
}

func (uh *UserHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
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

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	body := services.UserRequestParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, "Error deserializing request body", 500)
		return
	}

	createdUser, err := uh.userService.CreateUser(r.Context(), body)

	if err != nil {
		http.Error(w, "Error creating User", 500)
		fmt.Printf("Error creating user in db:\t%v", err.Error())
		return
	}

	uh.writeJSONResponse(w, createdUser, 200)
}

func (uh *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("id")
	userUUID, err := id.Parse(userId)
	if err != nil {
		http.Error(w, "Incorrectly formatted userID", http.StatusUnauthorized)
		fmt.Println("Error parsing requestID as UUID: ", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Erroneous Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error getting bearer token: ", err)
		return
	}

	err = uh.authService.ValidateJWT(token, userUUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonUpdateUser := services.UserRequestParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jsonUpdateUser)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		fmt.Println("Error decoding request body: ", err)
		return
	}

	updatedUser, err := uh.userService.UpdateUserById(r.Context(), jsonUpdateUser, token, userUUID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		fmt.Println("Experienced error while service tried to update user: ", err)
		return
	}

	uh.writeJSONResponse(w, updatedUser, http.StatusOK)
}

func (uh *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	userUUID, err := id.Parse(userId)
	if err != nil {
		http.Error(w, "Incorrectly formatted userID", http.StatusUnauthorized)
		fmt.Println("Error parsing requestID as UUID: ", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Erroneous Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error getting bearer token: ", err)
		return
	}

	deletedUser, err := uh.userService.DeleteUserByID(r.Context(), userUUID, token, uh.authService)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		fmt.Println("Error deleting user: ", err)
		return
	}

	uh.writeJSONResponse(w, deletedUser, http.StatusOK)
}
