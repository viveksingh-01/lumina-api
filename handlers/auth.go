package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/viveksingh-01/lumina-api/models"
	"github.com/viveksingh-01/lumina-api/utils"
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
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, utils.ErrorResponse{
			Error: "Invalid request method, please use POST method.",
		})
		return
	}
	// Validate request body to not be nil
	if r.Body == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Decoding the request body in JSON to struct format
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Error occurred while decoding request JSON", err.Error())
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Check if the user already exists based on the 'userId'
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": user.UserID}).Decode(&user)
	if err == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Username already exists, please try with a different one.",
		})
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Println("Database error: " + err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "An internal error occurred, please try again.",
		})
		return
	}

	// Generate hashed-password and store as password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Error occurred while hashing password:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	if _, err := userCollection.InsertOne(context.TODO(), user); err != nil {
		log.Println("Error occurred while inserting user's record to DB:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}
	log.Printf("New user registered: %s", user.UserID)

	// Write the response back to the client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful. You can now log in.",
	})
}
