package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/viveksingh-01/lumina-api/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(c *mongo.Collection) {
	userCollection = c
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
	http.SetCookie(w, createCookie(token, 0))
}

func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, createCookie("", -1))
}

func createCookie(value string, maxAge int) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	return cookie
}
