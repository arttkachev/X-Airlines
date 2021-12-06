package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var flightCollection *mongo.Collection

func GetFlightRepository() mongo.Collection {
	return *flightCollection
}

func SetFlightRepository(newCollection *mongo.Collection) {
	flightCollection = newCollection
}
