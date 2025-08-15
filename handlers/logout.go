package handlers

import (
	"encoding/json"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Logged out successfully!",
	})
}
