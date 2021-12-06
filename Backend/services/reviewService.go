package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection

func GetReviewRepository() mongo.Collection {
	return *reviewCollection
}

func SetReviewRepository(newCollection *mongo.Collection) {
	reviewCollection = newCollection
}
