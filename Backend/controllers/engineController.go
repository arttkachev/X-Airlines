package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateEngine(c *gin.Context) {
	var engine aircraft.Engine
	err := c.ShouldBindJSON(&engine)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	engine.ID = primitive.NewObjectID()
	engineService := services.GetEngineService()
	collection := engineService.Collection
	_, err = collection.InsertOne(ctx, engine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove engine data from Redis")
	engineService.RedisClient.Del("engines")
	c.JSON(http.StatusOK, engine)
}

func GetEngines(c *gin.Context) {
	var engineService = services.GetEngineService()
	val, err := engineService.RedisClient.Get("engines").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cur, err := engineService.Collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		defer cur.Close(ctx)
		engines := make([]aircraft.Engine, 0)
		for cur.Next(ctx) {
			var engine aircraft.Engine
			decodeErr := cur.Decode(&engine)
			if decodeErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": decodeErr.Error()})
				return
			}
			engines = append(engines, engine)
		}
		data, _ := json.Marshal(engines)
		engineService.RedisClient.Set("engines", string(data), 0)
		c.JSON(http.StatusOK, engines)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		engines := make([]aircraft.Engine, 0)
		json.Unmarshal([]byte(val), &engines)
		c.JSON(http.StatusOK, engines)
	}
}

func GetEngineById(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	var engineService = services.GetEngineService()
	val, err := engineService.RedisClient.Get("engines/" + id).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var engine aircraft.Engine
		err = engineService.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&engine)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		if engine.ID == primitive.NilObjectID {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No such engine"})
			return
		}

		data, _ := json.Marshal(engine)
		engineService.RedisClient.Set("engines/"+id, string(data), 0)
		c.JSON(http.StatusOK, engine)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		var engine aircraft.Engine
		json.Unmarshal([]byte(val), &engine)
		c.JSON(http.StatusOK, engine)
	}
}

func UpdateEngine(c *gin.Context) {
	var engine aircraft.Engine
	err := c.ShouldBindJSON(&engine)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	engineService := services.GetEngineService()
	collection := engineService.Collection
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"owningAircraft", bson.D{
			{"$cond", bson.D{
				{"if", engine.OwningAircraft != primitive.NilObjectID},
				{"then", engine.OwningAircraft},
				{"else", "$owningAircraft"}}}}},
		{"model", bson.D{
			{"$cond", bson.D{
				{"if", engine.Model != ""},
				{"then", engine.Model},
				{"else", "$model"}}}}},

		{"totalTime", bson.D{
			{"$cond", bson.D{
				{"if", engine.TotalTime != nil},
				{"then", engine.TotalTime},
				{"else", "$totalTime"}}}}},

		{"tbo", bson.D{
			{"$cond", bson.D{
				{"if", engine.TBO != nil},
				{"then", engine.TBO},
				{"else", "$tbo"}}}}},

		{"hst", bson.D{
			{"$cond", bson.D{
				{"if", engine.HST != nil},
				{"then", engine.HST},
				{"else", "$hst"}}}}}}}}
	_, err = collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove engine data from Redis")
	engineService.RedisClient.Del("engines")
	engineService.RedisClient.Del("engines/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The engine has been updated"})
	return
}

func DeleteEngine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	engineService := services.GetEngineService()
	collection := engineService.Collection
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an engine"})
		return
	}
	// clear cache
	log.Println("Remove engine data from Redis")
	engineService.RedisClient.Del("engines")
	engineService.RedisClient.Del("engines/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "An engine has been deleted"})
	return
}
