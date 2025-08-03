package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/services"
)

type UserHandler struct {
	db *database.Queries
}

func NewUserHandler(db *database.Queries) UserHandler {
	uh := UserHandler{db: db}
	return uh
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	type jsonUser struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		// CreatedAt  *time.Time `json:"createdAt"`
		// ModifiedAt *time.Time `json:"modifiedAt"`
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

	createdUser, err := uh.db.CreateUser(r.Context(), database.CreateUserParams{
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
