package services

import (
	//"github.com/arttkachev/X-Airlines/Backend/api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// temp storage
//var users []models.User = make([]models.User, 0)
var collection *mongo.Collection

func GetRepository() mongo.Collection {
	return *collection
}

func SetRepository(newCollection *mongo.Collection) {
	collection = newCollection
}
