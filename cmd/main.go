package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Lumina API.")

	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}
}
