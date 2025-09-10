package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	UserCollection *mongo.Collection
)

// SetCollections initializes all database collections
func SetCollections(db *mongo.Database) {
	UserCollection = db.Collection("users")
}
