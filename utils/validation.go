package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ValidateRequestMethod validates that the request method is POST
func ValidateRequestMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "Invalid request method, please use POST method.",
		})
		return false
	}
	return true
}

// ValidateRequestBody validates that the request body is not nil
func ValidateRequestBody(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil {
		SendErrorResponse(w, http.StatusBadRequest, ErrorResponse{
			Error: "Request body cannot be empty",
		})
		defer r.Body.Close()
		return false
	}
	return true
}

// DecodeToJSON decodes the request body in JSON to struct format
func DecodeToJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Println("Error occurred while decoding request JSON", err.Error())
		SendErrorResponse(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return false
	}
	return true
}
