package config

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

var Client *genai.Client

func InitializeGemini() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	var err error
	Client, err = genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Printf("Failed to initialize Gemini client: %v", err)
	}
}
