package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	aircraft "github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	"github.com/arttkachev/X-Airlines/Backend/api/models/airline"
	trackerdata "github.com/arttkachev/X-Airlines/Backend/api/models/trackerData"
	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAircraft(c *gin.Context) {
	var aircraft aircraft.Aircraft
	err := c.ShouldBindJSON(&aircraft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraft.ID = primitive.NewObjectID()
	aircraftService := services.GetAircraftService()
	if aircraft.Engines == nil {
		aircraft.Engines = make([]primitive.ObjectID, 0)
	}
	if aircraft.General.History == nil {
		aircraft.General.History = make([]primitive.ObjectID, 0)
	}
	if aircraft.Tags == nil {
		aircraft.Tags = make([]string, 0)
	}
	if aircraft.TrackerData == nil {
		aircraft.TrackerData = &trackerdata.TrackerData{FlightHistory: make([]primitive.ObjectID, 0)}
	}
	_, err = aircraftService.Collection.InsertOne(ctx, aircraft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	c.JSON(http.StatusOK, aircraft)
}

func GetAircraft(c *gin.Context) {
	aircraftService := services.GetAircraftService()
	airplanes := make([]aircraft.Aircraft, 0)
	val, err := aircraftService.RedisClient.Get("aircraft").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cur, err := aircraftService.Collection.Find(ctx, bson.M{})
		defer cur.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		for cur.Next(ctx) {
			var airplane aircraft.Aircraft
			err = cur.Decode(&airplane)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			airplanes = append(airplanes, airplane)
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(airplanes)
		aircraftService.RedisClient.Set("aircraft", string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airplanes)
	}
	c.JSON(http.StatusOK, airplanes)
}

func GetAircraftById(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	var aircraftService = services.GetAircraftService()
	var airplane aircraft.Aircraft
	val, err := aircraftService.RedisClient.Get("aircraft/" + id).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&airplane)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		if airplane.ID == primitive.NilObjectID {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No such aircraft"})
			return
		}
		data, _ := json.Marshal(airplane)
		aircraftService.RedisClient.Set("aircraft/"+id, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airplane)
	}
	c.JSON(http.StatusOK, airplane)
}

func GetAircraftByType(c *gin.Context) {
	airplaneQuery := c.Query("aircraft")
	aircraftService := services.GetAircraftService()
	foundAirplanes := make([]aircraft.Aircraft, 0)
	val, err := aircraftService.RedisClient.Get("aircraft/" + airplaneQuery).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cur, err := aircraftService.Collection.Find(ctx, bson.M{"general.name": airplaneQuery})
		defer cur.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		for cur.Next(ctx) {
			var airplane aircraft.Aircraft
			err = cur.Decode(&airplane)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			foundAirplanes = append(foundAirplanes, airplane)
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(foundAirplanes)
		aircraftService.RedisClient.Set("aircraft/"+airplaneQuery, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &foundAirplanes)
	}
	if len(foundAirplanes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Aircraft not found"})
		return
	}
	c.JSON(http.StatusOK, foundAirplanes)
}

func DeleteAircraft(c *gin.Context) {
	aircraftService := services.GetAircraftService()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	deleteResult, _ := aircraftService.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on deleting an aircraft"})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "An aircraft has been deleted"})
}

func UpdateAirframe(c *gin.Context) {
	var airframe aircraft.Airframe
	err := c.ShouldBindJSON(&airframe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"airframe.totalTime", bson.D{
			{"$cond", bson.D{
				{"if", airframe.TotalTime != nil},
				{"then", airframe.TotalTime},
				{"else", "$airframe.totalTime"}}}}},

		{"airframe.totalLandings", bson.D{
			{"$cond", bson.D{
				{"if", airframe.TotalLandings != nil},
				{"then", airframe.TotalLandings},
				{"else", "$airframe.totalLandings"}}}}},

		{"airframe.airframeNotes", bson.D{
			{"$cond", bson.D{
				{"if", airframe.AirframeNotes != ""},
				{"then", airframe.AirframeNotes},
				{"else", "$airframe.airframeNotes"}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft airframe has been updated"})
}

func UpdateExterior(c *gin.Context) {
	var exterior aircraft.Exterior
	err := c.ShouldBindJSON(&exterior)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"exterior.yearPainted", bson.D{
			{"$cond", bson.D{
				{"if", exterior.YearPainted != nil},
				{"then", exterior.YearPainted},
				{"else", "$exterior.yearPainted"}}}}},

		{"exterior.notes", bson.D{
			{"$cond", bson.D{
				{"if", exterior.Notes != ""},
				{"then", exterior.Notes},
				{"else", "$exterior.notes"}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft exterior has been updated"})
}

func UpdateInterior(c *gin.Context) {
	var interior aircraft.Interior
	err := c.ShouldBindJSON(&interior)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"interior.yearInterior", bson.D{
			{"$cond", bson.D{
				{"if", interior.YearInterior != nil},
				{"then", interior.YearInterior},
				{"else", "$interior.yearInterior"}}}}},

		{"interior.numberOfSeats", bson.D{
			{"$cond", bson.D{
				{"if", interior.NumberOfSeats != nil},
				{"then", interior.NumberOfSeats},
				{"else", "$interior.numberOfSeats"}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft interior has been updated"})
}

func UpdateCockpit(c *gin.Context) {
	var cockpit aircraft.Cockpit
	err := c.ShouldBindJSON(&cockpit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"cockpit.glassCockpit", bson.D{
			{"$cond", bson.D{
				{"if", cockpit.GlassCockpit != nil},
				{"then", cockpit.GlassCockpit},
				{"else", "$cockpit.glassCockpit"}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft cockpit has been updated"})
}

func UpdateGeneral(c *gin.Context) {
	var general aircraft.General
	err := c.ShouldBindJSON(&general)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"general.name", bson.D{
			{"$cond", bson.D{
				{"if", general.Name != ""},
				{"then", general.Name},
				{"else", "$general.name"}}}}},

		{"general.icon", bson.D{
			{"$cond", bson.D{
				{"if", general.Icon != ""},
				{"then", general.Icon},
				{"else", "$general.icon"}}}}},

		{"general.year", bson.D{
			{"$cond", bson.D{
				{"if", general.Year != nil},
				{"then", general.Year},
				{"else", "$general.year"}}}}},

		{"general.manufacturer", bson.D{
			{"$cond", bson.D{
				{"if", general.Manufacturer != ""},
				{"then", general.Manufacturer},
				{"else", "$general.manufacturer"}}}}},

		{"general.model", bson.D{
			{"$cond", bson.D{
				{"if", general.Model != ""},
				{"then", general.Model},
				{"else", "$general.model"}}}}},

		{"general.registration", bson.D{
			{"$cond", bson.D{
				{"if", general.Registration != ""},
				{"then", general.Registration},
				{"else", "$general.registration"}}}}},

		{"general.condition", bson.D{
			{"$cond", bson.D{
				{"if", general.Condition != ""},
				{"then", general.Condition},
				{"else", "$general.condition"}}}}},

		{"general.description", bson.D{
			{"$cond", bson.D{
				{"if", general.Description != ""},
				{"then", general.Description},
				{"else", "$general.description"}}}}},

		{"general.location", bson.D{
			{"$cond", bson.D{
				{"if", general.Location != ""},
				{"then", general.Location},
				{"else", "$general.location"}}}}},

		{"general.isOperating", bson.D{
			{"$cond", bson.D{
				{"if", general.IsOperating != nil},
				{"then", general.IsOperating},
				{"else", "$general.isOperating"}}}}},

		{"general.history", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{general.History, bson.A{}}}}}}}, "$general.history"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$general.history", general.History}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$general.history", bson.D{{"$ifNull", bson.A{general.History, bson.A{}}}}}}}}}}}},

		{"general.price", bson.D{
			{"$cond", bson.D{
				{"if", general.Price != nil},
				{"then", general.Price},
				{"else", "$general.price"}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft general information has been updated"})
}

func UpdatePerformance(c *gin.Context) {
	var performance aircraft.Performance
	err := c.ShouldBindJSON(&performance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"performance.range", bson.D{
			{"$cond", bson.D{
				{"if", performance.Range != nil},
				{"then", performance.Range},
				{"else", "$performance.range"}}}}},

		{"performance.cruiseSpeed", bson.D{
			{"$cond", bson.D{
				{"if", performance.CruiseSpeed != nil},
				{"then", performance.CruiseSpeed},
				{"else", "$performance.cruiseSpeed"}}}}},

		{"performance.maxSpeed", bson.D{
			{"$cond", bson.D{
				{"if", performance.MaxSpeed != nil},
				{"then", performance.MaxSpeed},
				{"else", "$performance.maxSpeed"}}}}},

		{"performance.ceiling", bson.D{
			{"$cond", bson.D{
				{"if", performance.Ceiling != nil},
				{"then", performance.Ceiling},
				{"else", "$performance.ceiling"}}}}},

		{"performance.maxTakeoffWeight", bson.D{
			{"$cond", bson.D{
				{"if", performance.MaxTakeoffWeight != nil},
				{"then", performance.MaxTakeoffWeight},
				{"else", "$performance.maxTakeoffWeight"}}}}},

		{"performance.maxLandingWeight", bson.D{
			{"$cond", bson.D{
				{"if", performance.MaxLandingWeight != nil},
				{"then", performance.MaxLandingWeight},
				{"else", "$performance.maxLandingWeight"}}}}},

		{"performance.maxZeroFuelWeight", bson.D{
			{"$cond", bson.D{
				{"if", performance.MaxZeroFuelWeight != nil},
				{"then", performance.MaxZeroFuelWeight},
				{"else", "$performance.maxZeroFuelWeight"}}}}},

		{"performance.fuelCapacity", bson.D{
			{"$cond", bson.D{
				{"if", performance.FuelCapacity != nil},
				{"then", performance.FuelCapacity},
				{"else", "$performance.fuelCapacity"}}}}},

		{"performance.takeoffDistance", bson.D{
			{"$cond", bson.D{
				{"if", performance.TakeoffDistance != nil},
				{"then", performance.TakeoffDistance},
				{"else", "$performance.takeoffDistance"}}}}},

		{"performance.wingspan", bson.D{
			{"$cond", bson.D{
				{"if", performance.Wingspan != nil},
				{"then", performance.Wingspan},
				{"else", "$performance.wingspan"}}}}}}}}
	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft performance has been updated"})
}

func UpdateEngines(c *gin.Context) {
	var airplane aircraft.Aircraft
	aircraftService := services.GetAircraftService()
	engineService := services.GetEngineService()
	err := c.ShouldBindJSON(&airplane)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	aircraftObjectId, _ := primitive.ObjectIDFromHex(id)

	// update current owning aircraft info
	for _, x := range airplane.Engines {
		var engine aircraft.Engine
		engineId := x.Hex()
		engineObjectId, _ := primitive.ObjectIDFromHex(engineId)
		engineVal, err := engineService.RedisClient.Get("engines/" + engineId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = engineService.Collection.FindOne(ctx, bson.M{"_id": engineObjectId}).Decode(&engine)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
			data, _ := json.Marshal(engine)
			engineService.RedisClient.Set("engines/"+engineId, string(data), 0)

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(engineVal), &engine)
		}
		if engine.OwningAircraft != primitive.NilObjectID {
			formerOwningAircraftObjectId, _ := primitive.ObjectIDFromHex(engine.OwningAircraft.Hex())
			var formerOwningAircraft aircraft.Aircraft
			formerOwningAircraftVal, err := aircraftService.RedisClient.Get("aircraft/" + engine.OwningAircraft.Hex()).Result()
			if err == redis.Nil {
				log.Printf("Request to MongoDB")
				err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": formerOwningAircraftObjectId}).Decode(&formerOwningAircraft)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error()})
					return
				}
				data, _ := json.Marshal(formerOwningAircraft)
				aircraftService.RedisClient.Set("aircraft/"+engine.OwningAircraft.Hex(), string(data), 0)
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			} else {
				log.Printf("Request to Redis")
				json.Unmarshal([]byte(formerOwningAircraftVal), &formerOwningAircraft)
			}

			if aircraftObjectId != formerOwningAircraftObjectId {
				var formerAircraft aircraft.Aircraft
				formerAircraft.Engines = append(formerAircraft.Engines, x)
				filter := bson.D{{"_id", formerOwningAircraftObjectId}}
				update := bson.D{{"$set", bson.D{
					{"engines", bson.D{
						{"$cond", bson.D{
							{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{formerAircraft.Engines, bson.A{}}}}}}}, formerOwningAircraft.Engines}}}},
							{"then", bson.D{{"$setDifference", bson.A{formerOwningAircraft.Engines, formerAircraft.Engines}}}},
							{"else", bson.D{{"$concatArrays", bson.A{formerOwningAircraft.Engines, bson.D{{"$ifNull", bson.A{formerAircraft.Engines, bson.A{}}}}}}}}}}}}}}}
				_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error()})
					return
				}
				log.Println("Remove aircraft id data from Redis")
				aircraftService.RedisClient.Del("aircraft/" + engine.OwningAircraft.Hex())
			}
		}
		filter := bson.D{{"_id", engineObjectId}}
		update := bson.D{{"$set", bson.D{
			{"owningAircraft", bson.D{
				{"$cond", bson.D{
					{"if", aircraftObjectId != engine.OwningAircraft},
					{"then", aircraftObjectId},
					{"else", primitive.NilObjectID}}}}}}}}
		_, err = engineService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		log.Println("Remove engine id data from Redis")
		engineService.RedisClient.Del("engines/" + engineId)
	}
	filter := bson.D{{"_id", aircraftObjectId}}
	update := bson.D{{"$set", bson.D{
		{"engines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{airplane.Engines, bson.A{}}}}}}}, "$engines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$engines", airplane.Engines}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$engines", bson.D{{"$ifNull", bson.A{airplane.Engines, bson.A{}}}}}}}}}}}}}}}
	// update engines
	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove engines data from Redis")
	engineService.RedisClient.Del("engines")
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft engines have been updated"})
}

func UpdateTags(c *gin.Context) {
	var aircraft aircraft.Aircraft
	err := c.ShouldBindJSON(&aircraft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"tags", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{aircraft.Tags, bson.A{}}}}}}}, "$tags"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$tags", aircraft.Tags}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$tags", bson.D{{"$ifNull", bson.A{aircraft.Tags, bson.A{}}}}}}}}}}}}}}}

	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft tags have been updated"})
}

func UpdateOwner(c *gin.Context) {
	var aircraft aircraft.Aircraft
	err := c.ShouldBindJSON(&aircraft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftService := services.GetAircraftService()
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"owner", bson.D{
			{"$cond", bson.D{
				{"if", aircraft.Owner != primitive.NilObjectID},
				{"then", aircraft.Owner},
				{"else", "$owner"}}}}}}}}
	_, err = aircraftService.Collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove aircraft data from Redis")
	aircraftService.RedisClient.Del("aircraft")
	aircraftService.RedisClient.Del("aircraft/" + id)
	c.JSON(http.StatusOK, gin.H{
		"message": "The owner has been updated"})
}

func GetOwnerData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftId := c.Param("id")
	var airplane aircraft.Aircraft
	aircraftService := services.GetAircraftService()
	val, err := aircraftService.RedisClient.Get("aircraft/" + aircraftId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		aircraftObjectId, err := primitive.ObjectIDFromHex(aircraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airplane)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(airplane)
		aircraftService.RedisClient.Set("aircraft/"+aircraftId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airplane)
	}

	if airplane.Owner == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Aircraft does not have an owner"})
		return
	}
	userId := airplane.Owner.Hex()
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

func GetEngineData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftId := c.Param("id")
	var airplane aircraft.Aircraft
	aircraftService := services.GetAircraftService()
	val, err := aircraftService.RedisClient.Get("aircraft/" + aircraftId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		aircraftObjectId, err := primitive.ObjectIDFromHex(aircraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airplane)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(airplane)
		aircraftService.RedisClient.Set("aircraft/"+aircraftId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airplane)
	}
	engineService := services.GetEngineService()
	var engines []aircraft.Engine
	for _, x := range airplane.Engines {
		engineId := x.Hex()
		engineObjectId, err := primitive.ObjectIDFromHex(engineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		var engine aircraft.Engine
		engineVal, err := engineService.RedisClient.Get("engines/" + engineId).Result()
		if err == redis.Nil {
			log.Printf("Request to MongoDB")
			err = engineService.Collection.FindOne(ctx, bson.M{"_id": engineObjectId}).Decode(&engine)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
			data, _ := json.Marshal(engine)
			aircraftService.RedisClient.Set("engines/"+engineId, string(data), 0)
			engines = append(engines, engine)

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		} else {
			log.Printf("Request to Redis")
			json.Unmarshal([]byte(engineVal), &engine)
			engines = append(engines, engine)

		}
	}
	if len(engines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The aircraft does not have engines"})
		return
	}
	c.JSON(http.StatusOK, engines)
}

func GetAirlineData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftId := c.Param("id")
	var airplane aircraft.Aircraft
	aircraftService := services.GetAircraftService()
	val, err := aircraftService.RedisClient.Get("aircraft/" + aircraftId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		aircraftObjectId, err := primitive.ObjectIDFromHex(aircraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		err = aircraftService.Collection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airplane)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(airplane)
		aircraftService.RedisClient.Set("aircraft/"+aircraftId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(val), &airplane)
	}
	if len(airplane.General.History) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The aircraft does not have an operator"})
		return
	}

	var airlane airline.Airline
	airlineService := services.GetAirlineService()
	airlineId := airplane.General.History[len(airplane.General.History)-1].Hex()
	airlineObjectId, err := primitive.ObjectIDFromHex(airlineId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	airlineVal, err := airlineService.RedisClient.Get("airlines/" + airlineId).Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		err = airlineService.Collection.FindOne(ctx, bson.M{"_id": airlineObjectId}).Decode(&airlane)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
			return
		}
		// Redis value has to be a string, so, we need to Marshal users first and put users on a Redis server
		data, _ := json.Marshal(airlane)
		aircraftService.RedisClient.Set("airlines/"+airlineId, string(data), 0)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		json.Unmarshal([]byte(airlineVal), &airlane)
	}
	c.JSON(http.StatusOK, airlane)
}
