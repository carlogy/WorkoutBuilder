package main

import (
	"fmt"
	"log"

	server "github.com/carlogy/WorkoutBuilder/internal/server"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	fmt.Println("Workout Builder!")

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)

	}

	// http.Handle("/", http.FileServer(http.Dir("./public")))

	// exercisesHandler := handlers.NewExerciseHandler()

	// http.HandleFunc("/api/exercises", exercisesHandler.GetExercises)

}
