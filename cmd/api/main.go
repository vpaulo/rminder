package main

import (
	"log"
	"rminder/internal/server"
)

func main() {
	server := server.NewServer("PORT")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %s", err)
	}

	log.Printf("Listening on %s...", server.Addr)
}
