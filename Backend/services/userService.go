package services

import (
	//"context"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var userHandler UserHandler

type UserHandler struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateUserHandler(collection *mongo.Collection, redisClient *redis.Client) *UserHandler {
	userHandler.Collection = collection
	userHandler.RedisClient = redisClient
	return &userHandler
}
func GetUserHandler() *UserHandler {
	return &userHandler
}
