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
	// Check if the user already exists based on the 'userId'
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

	// Set the token as cookie
	setCookie(w, token)

	// Send the response to client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "User registered successfully.",
		"data": map[string]string{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Validates request method
func validateRequestMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, utils.ErrorResponse{
			Error: "Invalid request method, please use POST method.",
		})
		return false
	}
	return true
}

// Validates request body to not be nil
func validateRequestBody(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Request body cannot be empty",
		})
		defer r.Body.Close()
		return false
	}
	return true
}

// Decodes the request body in JSON to struct format
func decodeToJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Println("Error occurred while decoding request JSON", err.Error())
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Invalid request body",
		})
		return false
	}
	return true
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
