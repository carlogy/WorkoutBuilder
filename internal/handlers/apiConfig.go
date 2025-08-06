package handlers

import "github.com/carlogy/WorkoutBuilder/internal/database"

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
