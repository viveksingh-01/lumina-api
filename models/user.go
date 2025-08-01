package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Email     string        `bson:"email"`
	Name      string        `bson:"name"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"createdAt"`
}
