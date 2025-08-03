package routes

import (
	"github.com/gorilla/mux"
	"github.com/viveksingh-01/lumina-api/handlers"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", handlers.Register)
	router.HandleFunc("/login", handlers.Login)
	router.HandleFunc("/api/chat", handlers.HandleChat).Methods("POST")
}
