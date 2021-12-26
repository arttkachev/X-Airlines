package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var airlineService AirlineService

type AirlineService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateAirlineService(collection *mongo.Collection, redisClient *redis.Client) *AirlineService {
	airlineService.Collection = collection
	airlineService.RedisClient = redisClient
	return &airlineService
}

func GetAirlineService() *AirlineService {
	return &airlineService
}
