package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/airline"
	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(c *gin.Context) {
	userService := services.GetUserService()
	// storage for found users
	users := make([]user.User, 0)
	val, err := userService.RedisClient.Get("users").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		// create context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// get a stream of documents (cursor) from mongo collection
		cur, err := userService.Collection.Find(ctx, bson.M{})
		// check on errors
		if err != nil {
			// return if error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// cursor must be closed on exit form function
		defer cur.Close(ctx)
		// loop through all users with mongo cursor
		for cur.Next(ctx) {
			var user user.User
			// decode mongo cursor into the user data type
			err = cur.Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			// append found user
			users = append(users, user)
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(users)
		userService.RedisClient.Set("users", string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &users)
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByAirline(c *gin.Context) {
	airline := c.Query("airlines")
	userService := services.GetUserService()
	users := make([]user.User, 0)
	val, err := userService.RedisClient.Get("users/" + airline).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cur, err := userService.Collection.Find(ctx, bson.M{"airlines": airline})
		defer cur.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		for cur.Next(ctx) {
			var user user.User
			err = cur.Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			users = append(users, user)
		}
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Users not found"})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(users)
		userService.RedisClient.Set("users/"+airline, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &users)
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
	var userService = services.GetUserService()
	collection := userService.Collection
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
				{"else", "$balance"}}}}}}}}

	// update
	_, err = collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove user data from Redis")
	userService.RedisClient.Del("users")
	userService.RedisClient.Del("users/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been updated"})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	userService := services.GetUserService()
	var user user.User
	userVal, err := userService.RedisClient.Get("users/" + id).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		err = userService.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(userVal), &user)

	}
	airlineService := services.GetAirlineService()
	for _, x := range user.Airlines {
		airlineId := x.Hex()
		airlineObjectId, err := primitive.ObjectIDFromHex(airlineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		_, err = airlineService.Collection.DeleteOne(ctx, bson.M{"_id": airlineObjectId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		log.Println("Remove airline id data from Redis")
		airlineService.RedisClient.Del("airlines/" + airlineId)
	}
	deleteResult, _ := userService.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting a user"})
		return
	}
	// clear cache
	log.Println("Remove user data from Redis")
	userService.RedisClient.Del("users")
	userService.RedisClient.Del("users/" + id)
	log.Println("Remove airlines data from Redis")
	airlineService.RedisClient.Del("airlines")
	c.JSON(http.StatusOK, gin.H{
		"message": "A user has been deleted"})
}

func GetUserAirlinesData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userId := c.Param("id")
	var user user.User
	userService := services.GetUserService()
	val, err := userService.RedisClient.Get("users/" + userId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		aircraftObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		err = userService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(user)
		userService.RedisClient.Set("users/"+userId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &user)
	}
	airlineService := services.GetAirlineService()
	var airlines []airline.Airline
	for _, x := range user.Airlines {
		airlineId := x.Hex()
		airlineObjectId, err := primitive.ObjectIDFromHex(airlineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		var airline airline.Airline
		engineVal, err := airlineService.RedisClient.Get("airlines/" + airlineId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = airlineService.Collection.FindOne(ctx, bson.M{"_id": airlineObjectId}).Decode(&airline)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
			data, _ := json.Marshal(airline)
			airlineService.RedisClient.Set("airlines/"+airlineId, string(data), 0)
			airlines = append(airlines, airline)

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(engineVal), &airline)
			airlines = append(airlines, airline)

		}
	}
	if len(airlines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The user does not own airlines"})
		return
	}
	c.JSON(http.StatusOK, airlines)
}

func UpdateUserAirlines(c *gin.Context) {
	var user user.User
	userService := services.GetUserService()
	airlineService := services.GetAirlineService()
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	userObjectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", userObjectId}}
	update := bson.D{{"$set", bson.D{
		{"airlines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{user.Airlines, bson.A{}}}}}}}, "$airlines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$airlines", user.Airlines}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$airlines", bson.D{{"$ifNull", bson.A{user.Airlines, bson.A{}}}}}}}}}}}}}}}

	_, err = userService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	for _, x := range user.Airlines {
		var airline airline.Airline
		airlineId := x.Hex()
		airlineObjectId, _ := primitive.ObjectIDFromHex(airlineId)
		airlineVal, err := airlineService.RedisClient.Get("airlines/" + airlineId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = airlineService.Collection.FindOne(ctx, bson.M{"_id": airlineObjectId}).Decode(&airline)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
			data, _ := json.Marshal(airline)
			airlineService.RedisClient.Set("airlines/"+airlineId, string(data), 0)

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(airlineVal), &airline)
		}
		filter := bson.D{{"_id", airlineObjectId}}
		update := bson.D{{"$set", bson.D{
			{"owner", bson.D{
				{"$cond", bson.D{
					{"if", airline.Owner == userObjectId},
					{"then", primitive.NilObjectID},
					{"else", userObjectId}}}}}}}}
		_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		log.Println("Remove airline id data from Redis")
		airlineService.RedisClient.Del("airlines/" + airlineId)
	}
	log.Println("Remove users data from Redis")
	userService.RedisClient.Del("users")
	userService.RedisClient.Del("users/" + id)
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	c.JSON(http.StatusOK, gin.H{
		"message": "The user airlines have been updated"})
}
