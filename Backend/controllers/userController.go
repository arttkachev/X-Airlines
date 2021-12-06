package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(c *gin.Context) {
	var user user.User
	bindErr := c.ShouldBindJSON(&user) // ShouldBindJSON marshals the incoming request body into a struct passed in as an argument (it's user in our case)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// generate a unique id for a user
	user.ID = primitive.NewObjectID()
	// get db collection
	collection := services.GetUserRepository()
	// add a new user
	if user.Airlines == nil {
		user.Airlines = make([]string, 0)
	}
	_, insertErr := collection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, user) // sends a response with httpStatusOK and a newly created user as a JSON
}

func GetUsers(c *gin.Context) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetUserRepository()
	// get a stream of documents (cursor) from mongo collection
	cur, findErr := collection.Find(ctx, bson.M{})
	// cursor must be closed on exit form function
	defer cur.Close(ctx)
	// check on errors
	if findErr != nil {
		// return if error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": findErr.Error()})
		return
	}
	// storage for found users
	users := make([]user.User, 0)

	// loop through all users with mongo cursor
	for cur.Next(ctx) {
		var user user.User
		// decode mongo cursor into the user data type
		decodeErr := cur.Decode(&user)
		if decodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": decodeErr.Error()})
			return
		}
		// append found user
		users = append(users, user)
	}
	// respond back with found users
	c.JSON(http.StatusOK, users)
}

func GetUserByAirline(c *gin.Context) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// fetch airline from user query
	airline := c.Query("airlines")
	// get db collection
	collection := services.GetUserRepository()
	// get a stream of documents (cursor) from mongo collection by query data
	cur, findErr := collection.Find(ctx, bson.M{"airlines": airline})
	// cursor must be closed on exit form function
	defer cur.Close(ctx)
	// check on errors
	if findErr != nil {
		// return if error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": findErr.Error()})
		return
	}
	// storage for found users
	users := make([]user.User, 0)
	// loop through all users with mongo cursor
	for cur.Next(ctx) {
		var user user.User
		// decode mongo cursor into the user data type
		decodeErr := cur.Decode(&user)
		if decodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": decodeErr.Error()})
			return
		}
		// append found user
		users = append(users, user)
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var user user.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetUserRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	// create a filter and mongo aggregation conds for PUT request
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"name", bson.D{
			{"$cond", bson.D{
				{"if", user.Name == ""},
				{"then", "$name"},
				{"else", user.Name}}}}},

		{"email", bson.D{
			{"$cond", bson.D{
				{"if", user.Email == ""},
				{"then", "$email"},
				{"else", user.Email}}}}},

		{"password", bson.D{
			{"$cond", bson.D{
				{"if", user.Password == ""},
				{"then", "$password"},
				{"else", user.Password}}}}},

		{"isAdmin", bson.D{
			{"$cond", bson.D{
				{"if", user.IsAdmin != nil},
				{"then", user.IsAdmin},
				{"else", "$isAdmin"}}}}},

		{"balance", bson.D{
			{"$cond", bson.D{
				{"if", user.Balance != nil},
				{"then", user.Balance},
				{"else", "$balance"}}}}},

		{"airlines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{user.Airlines, bson.A{}}}}}}}, "$airlines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$airlines", user.Airlines}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$airlines", bson.D{{"$ifNull", bson.A{user.Airlines, bson.A{}}}}}}}}}}}}}}}

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been updated"})
	return
}

func DeleteUser(c *gin.Context) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetUserRepository()
	// fetch the user id from the request URL
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	// delete by id
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	// respond an error if something's wrong (for example, a wrong id)
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting a user"})
		return
	}
	// respond a message to a client if everything is ok
	c.JSON(http.StatusOK, gin.H{
		"message": "A user has been deleted"})
	return
}
