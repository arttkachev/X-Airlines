package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var engineCollection *mongo.Collection

func GetEngineRepository() mongo.Collection {
	return *engineCollection
}

func SetEngineRepository(newCollection *mongo.Collection) {
	engineCollection = newCollection
}
