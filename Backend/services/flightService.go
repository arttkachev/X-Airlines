package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var flightService FlightService

type FlightService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateFlightService(collection *mongo.Collection, redisClient *redis.Client) *FlightService {
	flightService.Collection = collection
	flightService.RedisClient = redisClient
	return &flightService
}
func GetFlightService() *FlightService {
	return &flightService
}
