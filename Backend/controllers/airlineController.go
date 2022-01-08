package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	"github.com/arttkachev/X-Airlines/Backend/api/models/airline"
	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAirline(c *gin.Context) {
	var airlineData airline.Airline
	err := c.ShouldBindJSON(&airlineData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineData.ID = primitive.NewObjectID()
	airlineService := services.GetAirlineService()
	if airlineData.Fleet == nil {
		airlineData.Fleet = make([]primitive.ObjectID, 0)
	}
	if airlineData.Reviews == nil {
		airlineData.Reviews = make([]primitive.ObjectID, 0)
	}
	if airlineData.Routes == nil {
		airlineData.Routes = make([]primitive.ObjectID, 0)
	}
	_, err = airlineService.Collection.InsertOne(ctx, airlineData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove airlines data from Redis")
	airlineService.RedisClient.Del("airlines")
	c.JSON(http.StatusOK, airlineData)
}

func DeleteAirline(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	var airline airline.Airline
	airlineService := services.GetAirlineService()
	airlineVal, err := airlineService.RedisClient.Get("airlines/" + id).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		err = airlineService.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&airline)
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
		json.Unmarshal([]byte(airlineVal), &airline)
	}
	userService := services.GetUserService()
	ownerId := airline.Owner.Hex()
	ownerObjectId, _ := primitive.ObjectIDFromHex(ownerId)
	var selfAirline []primitive.ObjectID
	selfAirline = append(selfAirline, objectId)
	filter := bson.D{{"_id", ownerObjectId}}
	update := bson.D{{"$set", bson.D{
		{"airlines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{selfAirline, bson.A{}}}}}}}, "$airlines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$airlines", selfAirline}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$airlines", bson.D{{"$ifNull", bson.A{selfAirline, bson.A{}}}}}}}}}}}}}}}

	_, err = userService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove user data from Redis")
	userService.RedisClient.Del("users")
	userService.RedisClient.Del("users/" + ownerId)
	aircraftService := services.GetAircraftService()
	for _, x := range airline.Fleet {
		var airplane aircraft.Aircraft
		aircraftId := x.Hex()
		aircraftObjectId, err := primitive.ObjectIDFromHex(aircraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		aircraftVal, err := aircraftService.RedisClient.Get("aircraft/" + aircraftId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = airlineService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(airplane)
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
			json.Unmarshal([]byte(aircraftVal), &airplane)
		}
		if len(airplane.General.History) > 0 {
			var selfAirline []primitive.ObjectID
			selfAirline = append(selfAirline, objectId)
			filter := bson.D{{"_id", aircraftObjectId}}
			update := bson.D{{"$set", bson.D{
				{"general.history", bson.D{
					{"$cond", bson.D{
						{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{selfAirline, bson.A{}}}}}}}, "$general.history"}}}},
						{"then", bson.D{{"$setDifference", bson.A{"$general.history", selfAirline}}}},
						{"else", bson.D{{"$concatArrays", bson.A{"$general.history", bson.D{{"$ifNull", bson.A{selfAirline, bson.A{}}}}}}}}}}}}}}}

			_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			log.Println("Remove aircraft id data from Redis")
			aircraftService.RedisClient.Del("aircraft/" + aircraftId)
		}
	}
	deleteResult, _ := airlineService.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an airline"})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	c.JSON(http.StatusOK, gin.H{
		"message": "An airline has been deleted"})
}

func GetAirlines(c *gin.Context) {
	airlineService := services.GetAirlineService()
	airlines := make([]airline.Airline, 0)
	val, err := airlineService.RedisClient.Get("airlines").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cur, err := airlineService.Collection.Find(ctx, bson.M{})
		defer cur.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}

		for cur.Next(ctx) {
			var airline airline.Airline
			err = cur.Decode(&airline)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			airlines = append(airlines, airline)
		}
		// Redis value has to be a string, so, we need to Marshal data first and put it on a Redis server
		data, _ := json.Marshal(airlines)
		airlineService.RedisClient.Set("airlines", string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airlines)
	}
	c.JSON(http.StatusOK, airlines)
}

func UpdateAirlineGeneral(c *gin.Context) {
	var general airline.General
	err := c.ShouldBindJSON(&general)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	id := c.Param("id")
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

	_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline general information has been updated"})
}

func UpdateReviews(c *gin.Context) {
	var airline airline.Airline
	err := c.ShouldBindJSON(&airline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"reviews", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Reviews, bson.A{}}}}}}}, "$reviews"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$reviews", airline.Reviews}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$reviews", bson.D{{"$ifNull", bson.A{airline.Reviews, bson.A{}}}}}}}}}}}}}}}

	_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline reviews have been updated"})
}

func UpdateRoutes(c *gin.Context) {
	var airline airline.Airline
	err := c.ShouldBindJSON(&airline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"reviews", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Routes, bson.A{}}}}}}}, "$routes"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$routes", airline.Routes}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$routes", bson.D{{"$ifNull", bson.A{airline.Routes, bson.A{}}}}}}}}}}}}}}}

	_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline routes have been updated"})
}

func UpdateFleet(c *gin.Context) {
	var airline airline.Airline
	err := c.ShouldBindJSON(&airline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"fleet", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airline.Fleet, bson.A{}}}}}}}, "$fleet"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$fleet", airline.Fleet}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$fleet", bson.D{{"$ifNull", bson.A{airline.Fleet, bson.A{}}}}}}}}}}}}}}}

	_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	aircraftService := services.GetAircraftService()
	for _, x := range airline.Fleet {
		newAircraftObjectId, _ := primitive.ObjectIDFromHex(x.Hex())
		var newAircraft aircraft.Aircraft
		newAircraftVal, err := aircraftService.RedisClient.Get("aircraft/" + x.Hex()).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": newAircraftObjectId}).Decode(&newAircraft)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			data, _ := json.Marshal(newAircraft)
			aircraftService.RedisClient.Set("aircraft/"+x.Hex(), string(data), 0)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(newAircraftVal), &newAircraft)
		}
		var airlineHistory []primitive.ObjectID
		airlineHistory = append(airlineHistory, objectId)
		filter := bson.D{{"_id", newAircraftObjectId}}
		update := bson.D{{"$set", bson.D{
			{"general.history", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airlineHistory, bson.A{}}}}}}}, "$general.history"}}}},
					{"then", bson.D{{"$setDifference", bson.A{"$general.history", airlineHistory}}}},
					{"else", bson.D{{"$concatArrays", bson.A{"$general.history", bson.D{{"$ifNull", bson.A{airlineHistory, bson.A{}}}}}}}}}}}}}}}
		_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		log.Println("Remove aircraft id data from Redis")
		aircraftService.RedisClient.Del("aircraft/" + x.Hex())
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	c.JSON(http.StatusOK, gin.H{
		"message": "The airline fleet has been updated"})
}

func GetFleetData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	airlineId := c.Param("id")
	var airline airline.Airline
	val, err := airlineService.RedisClient.Get("airlines/" + airlineId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		airlineObjectId, err := primitive.ObjectIDFromHex(airlineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
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
		json.Unmarshal([]byte(val), &airline)
	}
	var aircraftArray []aircraft.Aircraft
	for _, x := range airline.Fleet {
		aircraftId := x.Hex()
		aircraftObjectId, err := primitive.ObjectIDFromHex(aircraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		aircraftService := services.GetAircraftService()
		var aircraft aircraft.Aircraft
		aircraftVal, err := aircraftService.RedisClient.Get("aircraft/" + aircraftId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&aircraft)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			//Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
			data, _ := json.Marshal(aircraft)
			aircraftService.RedisClient.Set("aircraft/"+aircraftId, string(data), 0)
			aircraftArray = append(aircraftArray, aircraft)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(aircraftVal), &aircraft)
			aircraftArray = append(aircraftArray, aircraft)
		}
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
	err := c.ShouldBindJSON(&airline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"owner", bson.D{
			{"$cond", bson.D{
				{"if", airline.Owner != primitive.NilObjectID},
				{"then", airline.Owner},
				{"else", "$owner"}}}}}}}}

	_, err = airlineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove airline data from Redis")
	airlineService.RedisClient.Del("airlines")
	airlineService.RedisClient.Del("airlines/" + id)
	userService := services.GetUserService()
	userId := airline.Owner.Hex()
	var newAirlines []primitive.ObjectID
	newAirlines = append(newAirlines, objectId)
	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	userFilter := bson.D{{"_id", userObjectId}}
	userUpdate := bson.D{{"$set", bson.D{
		{"airlines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{newAirlines, bson.A{}}}}}}}, "$airlines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$airlines", newAirlines}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$airlines", bson.D{{"$ifNull", bson.A{newAirlines, bson.A{}}}}}}}}}}}}}}}

	_, err = userService.Collection.UpdateOne(ctx, userFilter, mongo.Pipeline{userUpdate})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove user data from Redis")
	userService.RedisClient.Del("users")
	userService.RedisClient.Del("users/" + userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "The owner has been updated"})
	return
}

func GetAirlineOwnerData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	airlineService := services.GetAirlineService()
	airlineId := c.Param("id")
	var airline airline.Airline
	val, err := airlineService.RedisClient.Get("airlines/" + airlineId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		airlineObjectId, err := primitive.ObjectIDFromHex(airlineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
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
		json.Unmarshal([]byte(val), &airline)
	}

	if airline.Owner == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Airline does not have an owner"})
		return
	}
	userId := airline.Owner.Hex()
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	var owner user.User
	userService := services.GetUserService()
	userVal, err := userService.RedisClient.Get("users/" + userId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		err = userService.Collection.FindOne(ctx, bson.M{"_id": userObjectId}).Decode(&owner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		//Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(owner)
		userService.RedisClient.Set("users/"+userId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(userVal), &owner)
	}
	if owner.ID == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No owner found"})
		return
	}
	c.JSON(http.StatusOK, owner)
}
