package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/services"
	"github.com/google/uuid"
)

type UserHandler struct {
	conf ApiConfig
}

type jsonUser struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

func NewUserHandler(c *ApiConfig) UserHandler {
	uh := UserHandler{conf: *c}
	return uh
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

	body := jsonUser{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, "Error deserializing request body", 500)
		return
	}

	hashPW, err := auth.HashPassword(body.Password)
	if err != nil {
		http.Error(w, "Error creating user", 500)
		fmt.Printf("Error hashing password:\t%v", err.Error())
		return
	}

	createdUser, err := uh.conf.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: services.NoneNullToNullString(body.FirstName),
		LastName:  services.NoneNullToNullString(body.LastName),
		Email:     body.Email,
		Password:  hashPW,
	})

	if err != nil {
		http.Error(w, "Error creating User", 500)
		fmt.Printf("Error creating user in db:\t%v", err.Error())
		return
	}

	u := services.ConvertDBUserToUser(createdUser)

	uh.writeJSONResponse(w, u, 200)
}

func (uh *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("id")
	userUUID, err := uuid.Parse(userId)
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

	tokenUserId, err := auth.ValidateJWT(token, uh.conf.secret)
	if err != nil {
		http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
		fmt.Println("Error validating token: ", err)
		return
	}

	if tokenUserId != userUUID {
		http.Error(w, "Invalid userId provided for update", http.StatusUnauthorized)
		fmt.Println("Token user and userID do not match: ", err)
		return
	}

	jsonUpdateUser := jsonUser{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jsonUpdateUser)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		fmt.Println("Error decoding request body: ", err)
		return
	}

	hashPW, err := auth.HashPassword(jsonUpdateUser.Password)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		fmt.Println("Error hashing password: ", err)
		return
	}

	updatedUser, err := uh.conf.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		FirstName: services.NoneNullToNullString(jsonUpdateUser.FirstName),
		LastName:  services.NoneNullToNullString(jsonUpdateUser.LastName),
		Password:  hashPW,
		ID:        userUUID,
	})
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		fmt.Println("Error updating user, ", err)
		return
	}

	convertedUser := services.ConvertDBUpdateUserToUser(updatedUser)

	uh.writeJSONResponse(w, convertedUser, http.StatusOK)
}

func (uh *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	userUUID, err := uuid.Parse(userId)
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

	tokenUserId, err := auth.ValidateJWT(token, uh.conf.secret)
	if err != nil {
		http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
		fmt.Println("Error validating token: ", err)
		return
	}

	if tokenUserId != userUUID {
		http.Error(w, "Invalid userId provided for update", http.StatusUnauthorized)
		fmt.Println("Token user and userID do not match: ", err)
		return
	}

	deletedUser, err := uh.conf.db.DeleteUserById(r.Context(), userUUID)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		fmt.Println("Error deleting user: ", err)
		return
	}

	jsonReponseUser := services.ConvertDBDeleteUserToUser(deletedUser)

	uh.writeJSONResponse(w, jsonReponseUser, http.StatusOK)
}
