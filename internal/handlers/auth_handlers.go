package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(as services.AuthService) AuthHandler {
	return AuthHandler{AuthService: &as}
}

func (ah *AuthHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
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

func (ah *AuthHandler) AuthenticateByEmail(w http.ResponseWriter, r *http.Request) {

	ea := services.EmailAuthRequestParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ea)
	if err != nil {
		http.Error(w, "Error deserializing request body", 500)
		fmt.Printf("Error deserializing request body:\t%v\n", err.Error())
		return
	}

	u, err := ah.AuthService.AuthenticateByEmail(r.Context(), ea)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			fmt.Println("Invalid password provided: ", err)
			return
		}

		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		fmt.Println("Error storing refresh token: ", err)
		return
	}

	ah.writeJSONResponse(w, u, http.StatusOK)
}

func (ah *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error getting bearer refresh token")
		return
	}

	updatedToken, err := ah.AuthService.RefreshToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error updating token in db: ", err)
		return
	}

	ah.writeJSONResponse(w, updatedToken, http.StatusOK)
}

func (ah *AuthHandler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error getting bearer refresh token: ", err)
		return
	}

	isRevoked, err := ah.AuthService.RevokeToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !isRevoked {
		http.Error(w, "Unexpected", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
