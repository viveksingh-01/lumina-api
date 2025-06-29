package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    string        `bson:"userId"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"createdAt"`
}
