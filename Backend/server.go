package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aircraftController "github.com/arttkachev/X-Airlines/Backend/controllers"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// Mongo connection
	// create a client
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")))
	if err != nil {
		log.Fatal(err)
	}

	// create a context (context is how long an OS is going to wait before a connection esablished)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // wait 10 seconds

	// connect to db
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// make sure to disconnect d when a main func exits. "defer" provides this possibility for us
	defer client.Disconnect(ctx)

	// check that cnnection works by printing a list of database names of the client
	database, err := client.ListDatabaseNames(ctx, bson.M{}) // params (context, filter for returned db namesgo )
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(database)
	services.SetUserRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USERS")))
	services.SetAircraftRepository(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("AIRCRAFT")))

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
	router.POST("/aircraft", aircraftController.CreateAircraft)
	router.DELETE("/aircraft/:id", aircraftController.DeleteAircraft)

	// listen and serve
	router.Run(port)
}
