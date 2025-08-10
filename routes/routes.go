package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/viveksingh-01/lumina-api/handlers"
	"github.com/viveksingh-01/lumina-api/middlewares"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", handlers.Register)
	router.HandleFunc("/login", handlers.Login)

	secured := router.PathPrefix("/api").Subrouter()
	// Use auth-middleware to secure paths with prefixed with '/api'
	secured.Use(middlewares.AuthMiddleware)

	secured.HandleFunc("/me", handlers.MeHandler).Methods(http.MethodGet)
	secured.HandleFunc("/chat", handlers.HandleChat)
}
