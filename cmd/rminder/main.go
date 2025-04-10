package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"

	"rminder/internal/login/authenticator"
	"rminder/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Port not found in .env file: %v", err)
	}

	log.Printf("Server listening on http://localhost:%d/", port)

	hostAndPort := fmt.Sprintf("0.0.0.0:%d", port)
	if err := http.ListenAndServe(hostAndPort, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
