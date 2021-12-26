package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var aircraftService AircraftService

type AircraftService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateAircraftService(collection *mongo.Collection, redisClient *redis.Client) *AircraftService {
	aircraftService.Collection = collection
	aircraftService.RedisClient = redisClient
	return &aircraftService
}
func GetAircraftService() *AircraftService {
	return &aircraftService
}
