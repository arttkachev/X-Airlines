package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aircraftController "github.com/arttkachev/X-Airlines/Backend/controllers"
	airlineController "github.com/arttkachev/X-Airlines/Backend/controllers"
	engineController "github.com/arttkachev/X-Airlines/Backend/controllers"
	userController "github.com/arttkachev/X-Airlines/Backend/controllers"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const port = ":3000"

// Endpoints
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Welcome to X-Airlines Backend"})
}

func main() {

	// load .env file
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// Mongo connection
	// create a client
	client, clientErr := mongo.NewClient(options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")))
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	// create a context (context is how long an OS is going to wait before a connection esablished)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // wait 10 seconds

	// connect to db
	connectionErr := client.Connect(ctx)
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}

	// make sure to disconnect d when a main func exits. "defer" provides this possibility for us
	defer client.Disconnect(ctx)

	// check that cnnection works by printing a list of database names of the client
	database, dbErr := client.ListDatabaseNames(ctx, bson.M{}) // params (context, filter for returned db namesgo )
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	fmt.Println(database)
	services.SetUserRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USERS")))
	services.SetAircraftRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("AIRCRAFT")))
	services.SetEngineRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("ENGINES")))
	services.SetAirlineRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("AIRLINES")))
	services.SetFlightRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("FLIGHTS")))
	services.SetReviewRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("REVIEWS")))
	services.SetRouteRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("ROUTES")))

	// Routing
	// create a router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// handlers
	// users
	router.GET("/", Welcome)
	router.GET("/users", userController.GetUsers)
	router.GET("/users/airline_filter", userController.GetUserByAirline)
	router.POST("/users", userController.CreateUser)
	router.PUT("/users/:id", userController.UpdateUser)
	router.DELETE("/users/:id", userController.DeleteUser)

	// aircraft
	router.GET("/aircraft", aircraftController.GetAircraft)
	router.GET("/aircraft/aircraft_filter", aircraftController.GetAircraftByType)
	router.POST("/aircraft", aircraftController.CreateAircraft)
	router.DELETE("/aircraft/:id", aircraftController.DeleteAircraft)
	router.PUT("/aircraft/:id/update_airframe", aircraftController.UpdateAirframe)
	router.PUT("/aircraft/:id/update_exterior", aircraftController.UpdateExterior)
	router.PUT("/aircraft/:id/update_interior", aircraftController.UpdateInterior)
	router.PUT("/aircraft/:id/update_engines", aircraftController.UpdateEngines)
	router.PUT("/aircraft/:id/update_cockpit", aircraftController.UpdateCockpit)
	router.PUT("/aircraft/:id/update_general", aircraftController.UpdateGeneral)
	router.PUT("/aircraft/:id/update_performance", aircraftController.UpdatePerformance)
	router.PUT("/aircraft/:id/update_tags", aircraftController.UpdateTags)
	router.PUT("/aircraft/:id/update_owner", aircraftController.UpdateOwner)
	router.GET("/aircraft/:id/get_owner", aircraftController.GetOwnerData)
	router.GET("/aircraft/:id/get_engines", aircraftController.GetEngineData)

	// engines
	router.GET("/engines", engineController.GetEngines)
	router.POST("/engines", engineController.CreateEngine)
	router.PUT("/engines/:id", engineController.UpdateEngine)
	router.DELETE("/engines/:id", engineController.DeleteEngine)

	// airlines
	router.GET("/airlines", airlineController.GetAirlines)
	router.POST("/airlines", airlineController.CreateAirline)
	router.DELETE("/airlines/:id", airlineController.DeleteAirline)

	router.PUT("/airlines/:id/update_general", airlineController.UpdateAirlineGeneral)
	router.PUT("/airlines/:id/update_review", airlineController.UpdateReviews)
	router.PUT("/airlines/:id/update_routes", airlineController.UpdateRoutes)
	router.PUT("/airlines/:id/update_fleet", airlineController.UpdateFleet)
	router.PUT("/airlines/:id/update_owner", airlineController.UpdateAirlineOwner)
	router.GET("/airlines/:id/get_fleet", airlineController.GetFleetData)
	router.GET("/airlines/:id/get_owner", airlineController.GetAirlineOwnerData)

	// listen and serve
	router.Run(port)
}
