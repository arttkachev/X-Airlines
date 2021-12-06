package aircraft

import "go.mongodb.org/mongo-driver/bson/primitive"

// aircraft general info model
type General struct {
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Icon         string               `json:"icon,omitempty" bson:"icon,omitempty"`
	Year         *uint16              `json:"year,omitempty" bson:"year,omitempty"`
	Manufacturer string               `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	Model        string               `json:"model,omitempty" bson:"model,omitempty"`
	Registration string               `json:"registration,omitempty" bson:"registration,omitempty"`
	Condition    string               `json:"condition,omitempty" bson:"condition,omitempty"`
	Description  string               `json:"description,omitempty" bson:"description,omitempty"`
	Location     string               `json:"location,omitempty" bson:"location,omitempty"`
	IsOperating  *bool                `json:"isOperating,omitempty" bson:"isOperating,omitempty"`
	History      []primitive.ObjectID `json:"history" bson:"history"`
	Price        *float32             `json:"price,omitempty" bson:"price,omitempty"`
}
