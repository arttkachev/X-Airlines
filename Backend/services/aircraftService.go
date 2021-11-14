package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var aircraftCollection *mongo.Collection

func GetAircraftRepository() mongo.Collection {
	return *aircraftCollection
}

func SetAircraftRepository(newCollection *mongo.Collection) {
	aircraftCollection = newCollection
}
