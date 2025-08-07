package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/services"
)

type UserHandler struct {
	conf ApiConfig
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
	type jsonUser struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
	}

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
		FirstName: services.NoneNullToNullString(*body.FirstName),
		LastName:  services.NoneNullToNullString(*body.LastName),
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

func (uh *UserHandler) AuthenticateByEmail(w http.ResponseWriter, r *http.Request) {

	type emailAuthentication struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	ea := emailAuthentication{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ea)
	if err != nil {
		http.Error(w, "Error deserializing request body", 500)
		fmt.Printf("Error deserializing request body:\t%v\n", err.Error())
		return
	}

	dbUser, err := uh.conf.db.GetUserByEmail(r.Context(), ea.Email)
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

	if ea.ExpiresInSeconds != nil {
		token, err := auth.MakeJWT(dbUser.ID, uh.conf.secret, time.Duration(*ea.ExpiresInSeconds*int(time.Second)))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Printf("Eror making JWT: %v", err)
			return
		}
		u := services.ConvertFullDBUserToUser(dbUser, &token)
		uh.writeJSONResponse(w, u, http.StatusOK)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, uh.conf.secret, time.Hour)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("Eror making JWT: %v", err)
		return
	}
	u := services.ConvertFullDBUserToUser(dbUser, &token)
	uh.writeJSONResponse(w, u, http.StatusOK)
}
