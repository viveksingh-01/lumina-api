package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, payload ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		Error:   payload.Error,
		Details: payload.Details,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"error":"Failed to encode error response"}`, http.StatusInternalServerError)
	}
}
