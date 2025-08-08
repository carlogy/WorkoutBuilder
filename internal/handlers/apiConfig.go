package handlers

import (
	"fmt"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	db     *database.Queries
	secret string
}

func NewApiConfig(db *database.Queries, secret string) *ApiConfig {
	return &ApiConfig{
		db:     db,
		secret: secret,
	}
}

func (ac *ApiConfig) ValidateJWTRequestHeader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			fmt.Println("Error getting token", token, err)
			return
		}

		id, err := auth.ValidateJWT(token, ac.secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			fmt.Println("Error validating token: ", id, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ac *ApiConfig) GetUserIDFromToken(headers http.Header) (uuid.UUID, error) {

	token, err := auth.GetBearerToken(headers)
	if err != nil {
		fmt.Println("Error getting token ", token, err)
		return uuid.Nil, err
	}

	id, err := auth.ValidateJWT(token, ac.secret)
	if err != nil {
		fmt.Println("Error validating token ", token, err)
		return uuid.Nil, err
	}
	return id, nil
}
