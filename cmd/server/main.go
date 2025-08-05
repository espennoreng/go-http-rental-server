package main

import (
	"log"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/inmemory"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

func main() {
	userRepo := inmemory.NewUserRepository()

	userService := services.NewUserService(userRepo)

	server := api.NewServer(userService)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}