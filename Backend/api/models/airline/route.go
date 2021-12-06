package airline

import (
	"github.com/arttkachev/X-Airlines/Backend/api/models/airport"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Route struct {
	ID      primitive.ObjectID   `json:"id" bson:"_id"`
	From    *airport.Airport     `json:"from,omitempty" bson:"from,omitempty"`
	To      *airport.Airport     `json:"to,omitempty" bson:"to,omitempty"`
	Flights []primitive.ObjectID `json:"flights" bson:"flights"`
}
