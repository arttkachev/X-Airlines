package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var routeService RouteService

type RouteService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateRouteService(collection *mongo.Collection, redisClient *redis.Client) *RouteService {
	routeService.Collection = collection
	routeService.RedisClient = redisClient
	return &routeService
}
func GetRouteService() *RouteService {
	return &routeService
}
