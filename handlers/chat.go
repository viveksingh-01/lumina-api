package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/viveksingh-01/lumina-api/models"
	"google.golang.org/genai"
)

var (
	Client *genai.Client
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Couldn't decode the JSON", http.StatusInternalServerError)
		return
	}

	log.Println("user:", req.UserID)
	log.Println("message:", req.Message)

	// TODO
}
