package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/viveksingh-01/lumina-api/config"
	"github.com/viveksingh-01/lumina-api/routes"
)

func main() {
	fmt.Println("Welcome to Lumina API.")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	config.ConnectToDB()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	config.InitializeGemini()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
