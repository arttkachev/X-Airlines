package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateEngine(c *gin.Context) {
	var engine aircraft.Engine
	bindErr := c.ShouldBindJSON(&engine)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	engine.ID = primitive.NewObjectID()
	collection := services.GetEngineRepository()
	_, insertErr := collection.InsertOne(ctx, engine)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, engine)
}

func GetEngines(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetEngineRepository()
	cur, findErr := collection.Find(ctx, bson.M{})
	defer cur.Close(ctx)
	if findErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": findErr.Error()})
		return
	}
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
	c.JSON(http.StatusOK, engines)
}

func GetEngineById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetEngineRepository()
	ObjectId, engineIdErr := primitive.ObjectIDFromHex(c.Param("id"))
	if engineIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": engineIdErr.Error()})
		return
	}
	var engine aircraft.Engine
	DecodeErr := collection.FindOne(ctx, bson.M{"_id": ObjectId}).Decode(&engine)
	if DecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": DecodeErr.Error()})
		return
	}
	if engine.ID == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No such engine"})
		return
	}
	c.JSON(http.StatusOK, engine)
}

func UpdateEngine(c *gin.Context) {
	var engine aircraft.Engine
	bindErr := c.ShouldBindJSON(&engine)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := services.GetEngineRepository()
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
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The engine has been updated"})
	return
}

func DeleteEngine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	collection := services.GetEngineRepository()
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an engine"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "An engine has been deleted"})
	return
}
