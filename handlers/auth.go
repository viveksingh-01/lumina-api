package handlers

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(c *mongo.Collection) {
	userCollection = c
}
