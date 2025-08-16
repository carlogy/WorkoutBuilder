package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	services "github.com/carlogy/WorkoutBuilder/internal/services"
)

type AuthHanlder struct {
	conf *ApiConfig
}

func NewAuthHandler(c *ApiConfig) AuthHanlder {
	return AuthHanlder{conf: c}
}

func (ah *AuthHanlder) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
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

func (ah *AuthHanlder) AuthenticateByEmail(w http.ResponseWriter, r *http.Request) {

	type emailAuthentication struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ea := emailAuthentication{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ea)
	if err != nil {
		http.Error(w, "Error deserializing request body", 500)
		fmt.Printf("Error deserializing request body:\t%v\n", err.Error())
		return
	}

	dbUser, err := ah.conf.db.GetUserByEmail(r.Context(), ea.Email)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("Error querying user from db by email:\t%v\n", err.Error())
		return
	}

	err = auth.CheckPasswordHash(ea.Password, dbUser.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("Error comparing pw to stored hash:\t%v\n", err.Error())
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, ah.conf.secret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("Eror making JWT: %v", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error creating refresh token: ", err)
		return
	}

	_, err = ah.conf.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		fmt.Println("Error storing refresh token: ", err)
		return
	}

	u := services.ConvertFullDBUserToUser(dbUser, &token, &refreshToken)

	ah.writeJSONResponse(w, u, http.StatusOK)
}

func (ah *AuthHanlder) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error getting bearer refresh token")
		return
	}

	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, "Error updating refresh token", http.StatusInternalServerError)
		fmt.Println("Error making refresh token: ", err)
		return
	}

	dbToken, err := ah.conf.db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:   newToken,
		Token_2: token,
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error updating token in db: ", err)
		return
	}

	type responseRefreshToken struct {
		RefreshToken string `json:"token"`
	}

	body := responseRefreshToken{RefreshToken: dbToken.Token}

	ah.writeJSONResponse(w, body, http.StatusOK)
}

func (ah *AuthHanlder) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error getting bearer refresh token: ", err)
		return
	}

	revokedToken, err := ah.conf.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Error revoking token in db: ", err)
		return
	}

	if !revokedToken.RevokedAt.Valid {
		http.Error(w, "Unexpected", http.StatusInternalServerError)
		fmt.Println("Token returned is not revoked:  ", revokedToken.RevokedAt)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
