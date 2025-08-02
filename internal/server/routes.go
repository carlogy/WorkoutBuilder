package server

import (
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/handlers"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("public")))

	eh := handlers.NewExerciseHandler(s.db)
	mux.HandleFunc("GET /api/exercises", eh.GetExercises)
	mux.HandleFunc("GET /api/exercises/{id}", eh.GetExerciseById)
	mux.HandleFunc("DELETE /api/exercises/{id}", eh.DeleteExerciseByID)
	mux.HandleFunc("POST /api/exercises", eh.CreateExercise)

	return mux
}
