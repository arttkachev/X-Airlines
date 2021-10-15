package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const port = ":3000"

// Endpoints
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Welcome to X-Airlines"})
}

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	connectionString := os.Getenv("CONNECTION_STRING")

	// Mongo connection
	// create a client
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
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

	// Routing
	// create a router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// handler of a main page
	router.GET("/", Welcome)

	// listen and serve
	router.Run(port)
}
