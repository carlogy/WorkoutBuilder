package handlers

import (
	"encoding/json"
	"fmt"
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

	writeJSONResponse(w, u, 200)
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

	// To do Make JWT and return

	if ea.ExpiresInSeconds == nil {
		token, err := auth.MakeJWT(dbUser.ID, uh.conf.secret, time.Duration(time.Hour))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Printf("Eror making JWT: %v", err)
		}
		u := services.ConvertFullDBUserToUser(dbUser, &token)
		writeJSONResponse(w, u, http.StatusOK)
	}

	token, err := auth.MakeJWT(dbUser.ID, uh.conf.secret, time.Duration(*ea.ExpiresInSeconds))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("Eror making JWT: %v", err)
	}
	u := services.ConvertFullDBUserToUser(dbUser, &token)

	writeJSONResponse(w, u, 200)
}
