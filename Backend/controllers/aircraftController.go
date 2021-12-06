package controllers

import (
	"context"
	"net/http"
	"time"

	aircraft "github.com/arttkachev/X-Airlines/Backend/api/models/aircraft"
	trackerdata "github.com/arttkachev/X-Airlines/Backend/api/models/trackerData"
	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAircraft(c *gin.Context) {
	// aircraft
	var aircraft aircraft.Aircraft
	bindErr := c.ShouldBindJSON(&aircraft)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErr.Error()})
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
	// first, create empty arrays for array fields in data model to avoid errors on PUT operations
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
	// insert one
	_, insertErr := collection.InsertOne(ctx, aircraft)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, aircraft) // sends a response with httpStatusOK and a newly created aircraft as a JSON
}

func GetAircraft(c *gin.Context) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
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
	// storage for found aircraft
	airplanes := make([]aircraft.Aircraft, 0)

	// loop through all aircraft with mongo cursor
	for cur.Next(ctx) {
		var airplane aircraft.Aircraft
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

func GetAircraftByType(c *gin.Context) {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// fetch aircraft from user query
	airplaneQuery := c.Query("aircraft")
	// get db collection
	collection := services.GetAircraftRepository()
	// get a stream of documents (cursor) from mongo collection by query data
	cur, findErr := collection.Find(ctx, bson.M{"general.name": airplaneQuery})
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
	foundAirplanes := make([]aircraft.Aircraft, 0)
	// loop through all users with mongo cursor
	for cur.Next(ctx) {
		var airplane aircraft.Aircraft
		// decode mongo cursor into the user data type
		decodeErr := cur.Decode(&airplane)
		if decodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": decodeErr.Error()})
			return
		}
		// append found user
		foundAirplanes = append(foundAirplanes, airplane)
	}
	if len(foundAirplanes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Aircraft not found"})
		return
	}
	c.JSON(http.StatusOK, foundAirplanes)
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

func UpdateAirframe(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var airframe aircraft.Airframe
	airframeBindErr := c.ShouldBindJSON(&airframe)
	if airframeBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": airframeBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	// create a filter and mongo aggregation conds for PUT request
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

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft airframe has been updated"})
	return
}

func UpdateExterior(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var exterior aircraft.Exterior
	exteriorBindErr := c.ShouldBindJSON(&exterior)
	if exteriorBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": exteriorBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// create a filter and mongo aggregation conds for PUT request
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

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft exterior has been updated"})
	return
}

func UpdateInterior(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var interior aircraft.Interior
	interiorBindErr := c.ShouldBindJSON(&interior)
	if interiorBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": interiorBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// create a filter and mongo aggregation conds for PUT request
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

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft interior has been updated"})
	return
}

func UpdateCockpit(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var cockpit aircraft.Cockpit
	cockpitBindErr := c.ShouldBindJSON(&cockpit)
	if cockpitBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": cockpitBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// create a filter and mongo aggregation conds for PUT request
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"cockpit.glassCockpit", bson.D{
			{"$cond", bson.D{
				{"if", cockpit.GlassCockpit != nil},
				{"then", cockpit.GlassCockpit},
				{"else", "$cockpit.glassCockpit"}}}}}}}}

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft cockpit has been updated"})
	return
}

func UpdateGeneral(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var general aircraft.General
	generalBindErr := c.ShouldBindJSON(&general)
	if generalBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": generalBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	// create a filter and mongo aggregation conds for PUT request
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

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft general information has been updated"})
	return
}

func UpdatePerformance(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var performance aircraft.Performance
	performanceBindErr := c.ShouldBindJSON(&performance)
	if performanceBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": performanceBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	// create a filter and mongo aggregation conds for PUT request
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
	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft performance has been updated"})
	return
}

func UpdateEngines(c *gin.Context) {
	var aircraft aircraft.Aircraft
	aircraftBindErr := c.ShouldBindJSON(&aircraft)
	if aircraftBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": aircraftBindErr.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := services.GetAircraftRepository()
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"engines", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{aircraft.Engines, bson.A{}}}}}}}, "$engines"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$engines", aircraft.Engines}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$engines", bson.D{{"$ifNull", bson.A{aircraft.Engines, bson.A{}}}}}}}}}}}}}}}

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft engines have been updated"})
	return
}

func UpdateTags(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var aircraft aircraft.Aircraft
	aircraftBindErr := c.ShouldBindJSON(&aircraft)
	if aircraftBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": aircraftBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// create a filter and mongo aggregation conds for PUT request
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"tags", bson.D{
			{"$cond", bson.D{
				{"if", bson.D{{"$in", bson.A{bson.D{{"$first", bson.A{bson.D{{"$ifNull", bson.A{aircraft.Tags, bson.A{}}}}}}}, "$tags"}}}},
				{"then", bson.D{{"$setDifference", bson.A{"$tags", aircraft.Tags}}}},
				{"else", bson.D{{"$concatArrays", bson.A{"$tags", bson.D{{"$ifNull", bson.A{aircraft.Tags, bson.A{}}}}}}}}}}}}}}}

	// update
	_, updateErr := collection.UpdateOne(ctx, filter, mongo.Pipeline{update})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "The aircraft tags have been updated"})
	return
}

func UpdateOwner(c *gin.Context) {
	// unmarshals the request body into a user var and check if no error occured
	var aircraft aircraft.Aircraft
	ownerBindErr := c.ShouldBindJSON(&aircraft)
	if ownerBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ownerBindErr.Error()})
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get db collection
	collection := services.GetAircraftRepository()
	// fetch "id" from the user input
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	// create a filter and mongo aggregation conds for PUT request
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"owner", bson.D{
			{"$cond", bson.D{
				{"if", aircraft.Owner != primitive.NilObjectID},
				{"then", aircraft.Owner},
				{"else", "$owner"}}}}}}}}

	// update
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

func GetOwnerData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftCollection := services.GetAircraftRepository()
	aircraftObjectId, aircraftIdErr := primitive.ObjectIDFromHex(c.Param("id"))
	if aircraftIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": aircraftIdErr.Error()})
		return
	}
	var airplane aircraft.Aircraft
	aircraftDecodeErr := aircraftCollection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airplane)
	if aircraftDecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": aircraftDecodeErr.Error()})
		return
	}
	if airplane.Owner == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Aircraft does not have an owner"})
		return
	}
	userObjectId, userIdErr := primitive.ObjectIDFromHex(airplane.Owner.Hex())
	if userIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": userIdErr.Error()})
		return
	}
	userCollection := services.GetUserRepository()
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

func GetEngineData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	aircraftCollection := services.GetAircraftRepository()
	aircraftObjectId, aircraftIdErr := primitive.ObjectIDFromHex(c.Param("id"))
	if aircraftIdErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": aircraftIdErr.Error()})
		return
	}
	var airplane aircraft.Aircraft
	aircraftDecodeErr := aircraftCollection.FindOne(ctx, bson.M{"_id": aircraftObjectId}).Decode(&airplane)
	if aircraftDecodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": aircraftDecodeErr.Error()})
		return
	}
	var engines []aircraft.Engine
	for _, x := range airplane.Engines {
		engineObjectId, engineIdErr := primitive.ObjectIDFromHex(x.Hex())
		if engineIdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": engineIdErr.Error()})
			return
		}
		engineCollection := services.GetEngineRepository()
		var engine aircraft.Engine
		engineDecodeErr := engineCollection.FindOne(ctx, bson.M{"_id": engineObjectId}).Decode(&engine)
		if engineDecodeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": engineDecodeErr.Error()})
			return
		}
		engines = append(engines, engine)
	}
	if len(engines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The aircraft does not have engines"})
		return
	}
	c.JSON(http.StatusOK, engines)
}
