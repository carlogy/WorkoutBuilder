package server

import (
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/handlers"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("public")))

	eh := handlers.NewExerciseHandler(s.ApiConfig)
	uh := handlers.NewUserHandler(s.ApiConfig)
	ueh := handlers.NewUserExerciseHanlder(s.ApiConfig)
	wh := handlers.NewWorkoutHandler(s.ApiConfig)

	mux.HandleFunc("GET /api/exercises", s.ValidateJWTRequestHeader(eh.GetExercises))
	mux.HandleFunc("GET /api/exercises/{id}", s.ValidateJWTRequestHeader(eh.GetExerciseById))
	mux.HandleFunc("DELETE /api/exercises/{id}", s.ValidateJWTRequestHeader(eh.DeleteExerciseByID))
	mux.HandleFunc("POST /api/exercises", s.ValidateJWTRequestHeader(eh.CreateExercise))

	mux.HandleFunc("POST /api/record-exercise", s.ValidateJWTRequestHeader(ueh.CreateUserExerciseHandler))
	mux.HandleFunc("GET /api/record-exercise/{id}", s.ValidateJWTRequestHeader(ueh.GetUserExerciseHandler))
	mux.HandleFunc("POST /api/record-exercise/{id}", s.ValidateJWTRequestHeader(ueh.UpdateUserExerciseHandler))
	mux.HandleFunc("DELETE /api/record-exercise/{id}", s.ValidateJWTRequestHeader(ueh.DeleteUserExerciseRecordById))

	mux.HandleFunc("POST /api/users", uh.CreateUser)
	mux.HandleFunc("POST /api/users/{id}", s.ValidateJWTRequestHeader(uh.UpdateUserById))
	mux.HandleFunc("DELETE /api/users/{id}", s.ValidateJWTRequestHeader(uh.DeleteUserById))

	mux.HandleFunc("POST /api/workouts", s.ValidateJWTRequestHeader(wh.CreateWorkoutHandler))

	mux.HandleFunc("POST /api/login", uh.AuthenticateByEmail)

	return mux
}
