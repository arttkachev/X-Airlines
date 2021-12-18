package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	"github.com/arttkachev/X-Airlines/Backend/api/models/airline"
	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAirline(c *gin.Context) {
	var airlineData airline.Airline
	bindErr := c.ShouldBindJSON(&airlineData)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineData.ID = primitive.NewObjectID()
	collection := services.GetAirlineRepository()
	if airlineData.Fleet == nil {
		airlineData.Fleet = make([]primitive.ObjectID, 0)
	}
	if airlineData.Reviews == nil {
		airlineData.Reviews = make([]primitive.ObjectID, 0)
	}
	if airlineData.Routes == nil {
		airlineData.Routes = make([]primitive.ObjectID, 0)
	}
	_, insertErr := collection.InsertOne(ctx, airlineData)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, airlineData)
}

func DeleteAirline(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	collection := services.GetAirlineRepository()
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an airline"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "An airline has been deleted"})
	return
}

func GetAirlines(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	cur, err := collection.Find(ctx, bson.M{})
	defer cur.Close(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	airlines := make([]airline.Airline, 0)
	for cur.Next(ctx) {
		var airline airline.Airline
		decodeErr := cur.Decode(&airline)
		if decodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": decodeErr.Error()})
			return
		}
		airlines = append(airlines, airline)
	}
	// respond back with found users
	c.JSON(http.StatusOK, airlines)
}

func UpdateAirlineGeneral(c *gin.Context) {
	var general airline.General
	generalBindErr := c.ShouldBindJSON(&general)
	if generalBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": generalBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"general.name", bson.D{
			{"$cond", bson.D{
				{"if", general.Name != ""},
				{"then", general.Name},
				{"else", "$general.name"}}}}},

		{"general.logo", bson.D{
			{"$cond", bson.D{
				{"if", general.Logo != ""},
				{"then", general.Logo},
				{"else", "$general.logo"}}}}},

		{"general.iata", bson.D{
			{"$cond", bson.D{
				{"if", general.IATA != ""},
				{"then", general.IATA},
				{"else", "$general.iata"}}}}},

		{"general.icao", bson.D{
			{"$cond", bson.D{
				{"if", general.ICAO != ""},
				{"then", general.ICAO},
				{"else", "$general.icao"}}}}},

		{"general.fleet", bson.D{
			{"$cond", bson.D{
				{"if", general.Fleet != nil},
				{"then", general.Fleet},
				{"else", "$general.fleet"}}}}},

		{"general.rating", bson.D{
			{"$cond", bson.D{
				{"if", general.Rating != nil},
				{"then", general.Rating},
				{"else", "$general.rating"}}}}}}}}

	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline general information has been updated"})
	return
}

func UpdateReviews(c *gin.Context) {
	var airline airline.Airline
	airlineBindErr := c.ShouldBindJSON(&airline)
	if airlineBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": airlineBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"reviews", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Reviews, bson.A{}}}}}}}, "$reviews"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$reviews", airline.Reviews}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$reviews", bson.D{{"$ifNull", bson.A{airline.Reviews, bson.A{}}}}}}}}}}}}}}}

	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline reviews have been updated"})
	return
}

func UpdateRoutes(c *gin.Context) {
	var airline airline.Airline
	airlineBindErr := c.ShouldBindJSON(&airline)
	if airlineBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": airlineBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"reviews", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Routes, bson.A{}}}}}}}, "$routes"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$routes", airline.Routes}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$routes", bson.D{{"$ifNull", bson.A{airline.Routes, bson.A{}}}}}}}}}}}}}}}

	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline routes have been updated"})
	return
}

func UpdateFleet(c *gin.Context) {
	var airline airline.Airline
	airlineBindErr := c.ShouldBindJSON(&airline)
	if airlineBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": airlineBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"fleet", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Fleet, bson.A{}}}}}}}, "$fleet"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$fleet", airline.Fleet}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$fleet", bson.D{{"$ifNull", bson.A{airline.Fleet, bson.A{}}}}}}}}}}}}}}}

	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline fleet has been updated"})
	return
}

func GetFleetData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineCollection := services.GetAirlineRepository()
	aircraftObjectId, aircraftIdErr := primitive.ObjectIDFromHex(c.Param("id"))
	if aircraftIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": aircraftIdErr.Error()})
		return
	}
	var airlineData airline.Airline
	airlineDecodeErr := airlineCollection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airlineData)
	if airlineDecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": airlineDecodeErr.Error()})
		return
	}
	var aircraftArray []aircraft.Aircraft
	for _, x := range airlineData.Fleet {
		fleetObjectId, fleetIdErr := primitive.ObjectIDFromHex(x.Hex())
		if fleetIdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fleetIdErr.Error()})
			return
		}
		aircraftCollection := services.GetAircraftRepository()
		var fleet aircraft.Aircraft
		fleetDecodeErr := aircraftCollection.FindOne(ctx, bson.M{"_id": fleetObjectId}).Decode(&fleet)
		if fleetDecodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fleetDecodeErr.Error()})
			return
		}
		aircraftArray = append(aircraftArray, fleet)
	}
	if len(aircraftArray) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The airline has no fleet"})
		return
	}
	c.JSON(http.StatusOK, aircraftArray)
}

func UpdateAirlineOwner(c *gin.Context) {
	var airline airline.Airline
	ownerBindErr := c.ShouldBindJSON(&airline)
	if ownerBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ownerBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAirlineRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"owner", bson.D{
			{"$cond", bson.D{
				{"if", airline.Owner != primitive.NilObjectID},
				{"then", airline.Owner},
				{"else", "$owner"}}}}}}}}

	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The owner has been updated"})
	return
}

func GetAirlineOwnerData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineCollection := services.GetAirlineRepository()
	airlineObjectId, airlineIdErr := primitive.ObjectIDFromHex(c.Param("id"))
	if airlineIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": airlineIdErr.Error()})
		return
	}
	var airline airline.Airline
	airlineDecodeErr := airlineCollection.FindOne(ctx, bson.M{"_id": airlineObjectId}).Decode(&airline)
	if airlineDecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": airlineDecodeErr.Error()})
		return
	}
	if airline.Owner == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Airline does not have an owner"})
		return
	}
	userObjectId, userIdErr := primitive.ObjectIDFromHex(airline.Owner.Hex())
	if userIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": userIdErr.Error()})
		return
	}
	var userService = services.GetUserService()
	userCollection := userService.Collection
	var owner user.User
	userDecodeErr := userCollection.FindOne(ctx, bson.M{"_id": userObjectId}).Decode(&owner)
	if userDecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": userDecodeErr.Error()})
		return
	}
	if owner.ID == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No owner found"})
		return
	}
	c.JSON(http.StatusOK, owner)
}
