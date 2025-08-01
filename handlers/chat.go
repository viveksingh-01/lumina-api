package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/viveksingh-01/lumina-api/models"
	"google.golang.org/genai"
)

var (
	Client   *genai.Client
	sessions = make(map[string]*genai.Chat)
	mu       sync.Mutex
)

const GEMINI_MODEL = "gemini-2.0-flash"

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

	mu.Lock()
	session, exists := sessions[req.UserID]
	if !exists {
		session, err := Client.Chats.Create(r.Context(), GEMINI_MODEL, nil, nil)
		if err != nil {
			log.Println(err)
		}
		sessions[req.UserID] = session
	}
	mu.Unlock()

	resp, err := session.SendMessage(r.Context(), genai.Part{Text: req.Message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.ChatResponse{Response: resp.Text()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
