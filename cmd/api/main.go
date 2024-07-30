package main

import (
	"fmt"
	"log"
	"rminder/internal/server"
)

func main() {

	server := server.NewServer()

	log.Printf("Listening on %s...", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
