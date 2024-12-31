package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"rminder/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer(key string) *http.Server {
	port, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		log.Fatalf("Port not found in .env file: %v", err)
	}

	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
