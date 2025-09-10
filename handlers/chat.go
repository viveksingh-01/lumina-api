package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/viveksingh-01/lumina-api/middlewares"
	"github.com/viveksingh-01/lumina-api/models"
	"github.com/viveksingh-01/lumina-api/utils"
	"google.golang.org/genai"
)

var (
	Client   *genai.Client
	sessions = make(map[string]*genai.Chat)
	mu       sync.Mutex
)

const GEMINI_MODEL = "gemini-2.0-flash"

func HandleChat(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	if !utils.ValidateRequestMethod(w, r) {
		return
	}
	if !utils.ValidateRequestBody(w, r) {
		return
	}

	var req models.ChatRequest
	if !utils.DecodeToJSON(w, r, &req) {
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
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ChatResponse{Response: resp.Text()})
}
