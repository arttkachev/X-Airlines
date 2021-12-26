package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var engineService EngineService

type EngineService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateEngineService(collection *mongo.Collection, redisClient *redis.Client) *EngineService {
	engineService.Collection = collection
	engineService.RedisClient = redisClient
	return &engineService
}
func GetEngineService() *EngineService {
	return &engineService
}
