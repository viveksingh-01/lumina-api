package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/viveksingh-01/lumina-api/models"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	// Check if the user already exists based on the 'userId'
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": user.UserID}).Decode(&user)
	if err == nil {
		http.Error(w, "Username already exists, please try with a different one.", http.StatusBadRequest)
		return
	}
	if err != mongo.ErrNoDocuments {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO
}
