package main

import (
	"log"
	"net/http"

	"example.com/auth_service/handlers"
)

func main() {
	//http.HandleFunc("/home", handlers.CreateUserHandler)
	http.HandleFunc("/create-token", handlers.GetTokensHandler)
	http.HandleFunc("/refresh-token", handlers.RefreshTokenHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
