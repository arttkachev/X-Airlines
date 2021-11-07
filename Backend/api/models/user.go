package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// user model
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	IsAdmin  *bool              `json:"isAdmin" bson:"isAdmin"`
	Balance  *int               `json:"balance" bson:"balance"`
	Airlines []string           `json:"airlines" bson:"airlines"`
}
