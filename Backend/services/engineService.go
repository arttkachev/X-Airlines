package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var engineService UserService

type EngineService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateEngineService(collection *mongo.Collection, redisClient *redis.Client) *UserService {
	engineService.Collection = collection
	engineService.RedisClient = redisClient
	return &userService
}
func GetEngineService() *UserService {
	return &userService
}
