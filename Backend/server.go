package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const port = ":3000"

var users []models.User

// Endpoints
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Welcome to X-Airlines"})
}

func AddUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user) // ShouldBindJSON marshals the incoming request body into a struct passed in as an argument (it's user in our case)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user.ID = xid.New().String() // xid lib set a unique ID
	users = append(users, user)
	c.JSON(http.StatusOK, user) // sends a response with httpStatusOK and a newly created user as a JSON
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users) //marshals the users array to JSON
}

func main() {

	users = make([]models.User, 0)

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

	// handlers
	router.GET("/", Welcome)
	router.POST("/users", AddUser)
	router.GET("/users", GetUsers)

	// listen and serve
	router.Run(port)
}
