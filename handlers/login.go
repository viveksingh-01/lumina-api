package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/viveksingh-01/lumina-api/models"
	"github.com/viveksingh-01/lumina-api/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if !validateRequestBody(w, r) {
		return
	}
	if !validateRequestMethod(w, r) {
		return
	}
	var req models.LoginRequest
	if !decodeToJSON(w, r, &req) {
		return
	}

	var user models.User
	// Check if the user exists
	err := userCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Println("User doesn't exist")
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "User doesn't exist, please create an account.",
		})
		return
	}
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "An internal error occurred, please try again.",
		})
		return
	}

	// Validate input password by comparing with hashed password stored in DB
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Invalid email or password, please try again.",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Println("Error generating JWT:", err.Error())
		utils.SendErrorResponse(w, http.StatusInternalServerError, utils.ErrorResponse{
			Error: "Couldn't process the request, please try again.",
		})
		return
	}
	setCookie(w, token)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "You've logged-in successfully!",
		"data": map[string]string{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
