package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func GetUserRepository() mongo.Collection {
	return *userCollection
}

func SetUserRepository(newCollection *mongo.Collection) {
	userCollection = newCollection
}
