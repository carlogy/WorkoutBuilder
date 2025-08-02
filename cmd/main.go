package main

import (
	"fmt"
	"log"
	"os"

	server "github.com/carlogy/WorkoutBuilder/internal/server"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env variables:\t%v", err)
	}

	fmt.Println("Workout Builder!")

	server := server.NewServer()

	err = server.ListenAndServe()
	log.Printf("Listening on Port:\t%v", os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Server failed: %v", err)

	}
}
