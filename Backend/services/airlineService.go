package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var airlineCollection *mongo.Collection

func GetAirlineRepository() mongo.Collection {
	return *airlineCollection
}

func SetAirlineRepository(newCollection *mongo.Collection) {
	airlineCollection = newCollection
}
