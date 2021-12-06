package trackerdata

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackerData struct {
	FlightHistory []primitive.ObjectID `json:"flightHistory" bson:"flightHistory"`
}
