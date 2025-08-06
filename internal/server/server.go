package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/handlers"

	_ "github.com/lib/pq"
)

type Server struct {
	port int
	*handlers.ApiConfig
}

func NewServer() *http.Server {

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	secret := os.Getenv("JWTSECRET")
	connStr := os.Getenv("WORKOUTBUILDER_DB_URL")
	if connStr == "" {
		log.Fatal("DB URL must be set")
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable open connection to database: %v", err)
	}

	dbQueries := database.New(dbConn)

	apiConfig := handlers.NewApiConfig(dbQueries, secret)

	NewServer := &Server{
		port:      port,
		ApiConfig: apiConfig,
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", NewServer.port),
		Handler:           NewServer.RegisterRoutes(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Minute,
	}

	return server
}
