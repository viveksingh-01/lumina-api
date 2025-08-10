package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/viveksingh-01/lumina-api/middlewares"
	"github.com/viveksingh-01/lumina-api/models"
	"github.com/viveksingh-01/lumina-api/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	objectID, err := bson.ObjectIDFromHex(claims.Subject)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, utils.ErrorResponse{
			Error: "Invalid user ID in token",
		})
		return
	}

	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data": map[string]string{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
