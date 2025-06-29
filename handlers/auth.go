package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/viveksingh-01/lumina-api/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(c *mongo.Collection) {
	userCollection = c
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method, please use POST method.", http.StatusMethodNotAllowed)
		return
	}
	// Validate request body to not be nil
	if r.Body == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decoding the request body in JSON to struct format
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error occurred while decoding request JSON", http.StatusBadRequest)
		return
	}

	// TODO
}
