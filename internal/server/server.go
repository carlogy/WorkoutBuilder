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

type Config struct {
	Port      int
	SecretKey string
	DBURI     string
	*database.Queries
}

type Server struct {
	*Config
	*handlers.ApiConfig
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	secret := os.Getenv("JWTSECRET")
	dbURI := os.Getenv("WORKOUTBUILDER_DB_URL")

	return &Config{
		Port:      port,
		SecretKey: secret,
		DBURI:     dbURI,
	}
}

func NewServer() *http.Server {

	cfg := NewConfig()

	if cfg.DBURI == "" {
		log.Fatal("DB URL must be set")
	}

	dbConn, err := sql.Open("postgres", cfg.DBURI)
	if err != nil {
		log.Fatalf("Unable open connection to database: %v", err)
	}

	dbQueries := database.New(dbConn)
	cfg.Queries = dbQueries

	apiCfg := handlers.NewApiConfig(cfg.Queries, cfg.SecretKey)

	NewServer := &Server{
		Config:    cfg,
		ApiConfig: apiCfg,
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", NewServer.Config.Port),
		Handler:           NewServer.RegisterRoutes(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Minute,
	}

	return server
}
