package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var routeCollection *mongo.Collection

func GetRouteRepository() mongo.Collection {
	return *routeCollection
}

func SetRouteRepository(newCollection *mongo.Collection) {
	routeCollection = newCollection
}
