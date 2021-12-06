package airline

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	User    string             `json:"user,omitempty" bson:"user,omitempty"`
	Avatar  string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Comment string             `json:"comment,omitempty" bson:"comment,omitempty"`
	Rating  *uint8             `json:"rating,omitempty" bson:"rating,omitempty"`
}
