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

func Register(w http.ResponseWriter, r *http.Request) {
	if !validateRequestMethod(w, r) {
		return
	}
	if !validateRequestBody(w, r) {
		return
	}
	var req models.RegisterRequest
	if !decodeToJSON(w, r, &req) {
		return
	}

	var user models.User
	// Check if the user already exists based on the 'email'
	err := userCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
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
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Error occurred while hashing password:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}

	user.Email = req.Email
	user.Name = req.Name
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	if _, err := userCollection.InsertOne(context.TODO(), user); err != nil {
		log.Println("Error occurred while inserting user's record to DB:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}
	log.Printf("New user registered: %s", user.Email)

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Println("Error generating JWT:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}

	// Send the response to client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "User registered successfully.",
		"token":   token,
		"data": map[string]string{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
