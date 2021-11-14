package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/fleet"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAircraft(c *gin.Context) {
	var aircraft fleet.Aircraft
	err := c.ShouldBindJSON(&aircraft) // ShouldBindJSON marshals the incoming request body into a struct passed in as an argument (aircraft in our case)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// generate a unique id for an aircraft
	aircraft.ID = primitive.NewObjectID()
	// get db collection
	collection := services.GetAircraftRepository()

	// create a new aircraft
	_, err = collection.InsertOne(ctx, aircraft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aircraft) // sends a response with httpStatusOK and a newly created user as a JSON
}

func GetAircraft(c *gin.Context) {

	// get db collection
	collection := services.GetAircraftRepository()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get a stream of documents (cursor) from mongo collection
	cur, err := collection.Find(ctx, bson.M{})
	// cursor must be closed on exit form function
	defer cur.Close(ctx)

	// check on errors
	if err != nil {
		// return if error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// our storage for found users
	airplanes := make([]fleet.Aircraft, 0)

	// loop through all users with mongo cursor
	for cur.Next(ctx) {
		var airplane fleet.Aircraft
		// decode mongo cursor into the user data type
		decodeErr := cur.Decode(&airplane)
		if decodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": decodeErr.Error()})
			return
		}
		// append found user
		airplanes = append(airplanes, airplane)
	}
	// respond back with found users
	c.JSON(http.StatusOK, airplanes)
}

func DeleteAircraft(c *gin.Context) {

	// get db collection
	collection := services.GetAircraftRepository()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch the user id from the request URL
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// delete by id
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	// respond an error if something's wrong (for example, a wrong id)
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an aircraft"})
		return
	}
	// respond a message to a client if everything is ok
	c.JSON(http.StatusOK, gin.H{
		"message": "An aircraft has been deleted"})
	return
}
