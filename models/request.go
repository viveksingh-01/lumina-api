package models

type RegisterRequest struct {
	Email    string `bson:"email"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}
