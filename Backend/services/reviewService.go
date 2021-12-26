package services

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewService ReviewService

type ReviewService struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
}

func CreateReviewService(collection *mongo.Collection, redisClient *redis.Client) *ReviewService {
	reviewService.Collection = collection
	reviewService.RedisClient = redisClient
	return &reviewService
}
func GetReviewService() *ReviewService {
	return &reviewService
}
