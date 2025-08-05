package main

import (
	"log"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/api"
)

func main() {
	server := api.NewServer()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}