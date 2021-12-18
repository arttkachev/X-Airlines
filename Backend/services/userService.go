package services

import (
	//"context"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var userService UserService

type UserService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateUserService(collection *mongo.Collection, redisClient *redis.Client) *UserService {
	userService.Collection = collection
	userService.RedisClient = redisClient
	return &userService
}
func GetUserService() *UserService {
	return &userService
}
