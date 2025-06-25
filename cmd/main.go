package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/viveksingh-01/lumina-api/routes"
)

func main() {
	fmt.Println("Welcome to Lumina API.")

	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// Use port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at port:", port)

	// Start the HTTP server and listen at the port-8080
	log.Fatal(http.ListenAndServe(":"+port, r))
}
