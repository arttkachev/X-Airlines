package airport

import (
	"github.com/arttkachev/X-Airlines/Backend/api/models/flight"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Airport struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id"`
	ICAO       string               `json:"icao,omitempty" bson:"icao,omitempty"`
	IATA       string               `json:"iata,omitempty" bson:"iata,omitempty"`
	Weather    primitive.ObjectID   `json:"weather,omitempty" bson:"weather,omitempty"`
	Arrivals   []flight.Flight      `json:"arrivals" bson:"arrivals"`
	Departures []flight.Flight      `json:"departures" bson:"departures"`
	TopTraffic []primitive.ObjectID `json:"topTraffic" bson:"topTraffic"`
}
