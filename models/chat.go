package models

type ChatRequest struct {
	Message string `json:"message"`
	UserID  string `json:"userId"`
}
