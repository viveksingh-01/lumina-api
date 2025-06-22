package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Lumina API.")

	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	// Use port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at port:", port)

	// Start the HTTP server and listen at the port-8080
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
