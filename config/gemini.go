package config

import (
	"context"
	"log"
	"os"

	"github.com/viveksingh-01/lumina-api/handlers"
	"google.golang.org/genai"
)

func InitializeGemini() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	c, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Failed to initialize Gemini client: %v", err)
	}
	handlers.Client = c
}
