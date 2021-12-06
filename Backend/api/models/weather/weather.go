package weather

import "go.mongodb.org/mongo-driver/bson/primitive"

type Weather struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Condition   string             `json:"condition,omitempty" bson:"condition,omitempty"`
	Temperature int8               `json:"temperature,omitempty" bson:"temperature,omitempty"`
}
