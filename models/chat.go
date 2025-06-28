package models

type ChatRequest struct {
	Message string `json:"message"`
	UserID  string `json:"userId"`
}

type ChatResponse struct {
	Response string `json:"response"`
}
