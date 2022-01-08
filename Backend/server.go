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
	auth "github.com/arttkachev/X-Airlines/Backend/services/auth"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	// init redis store for user session cookies
	store, _ := redisSession.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	// init Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// Ping Redis clinet
	redisStatus := redisClient.Ping()
	// print redist status result
	fmt.Println(redisStatus)

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
	AuthService := auth.AuthService{}
	services.CreateUserService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USERS")), redisClient)
	services.CreateAircraftService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("AIRCRAFT")), redisClient)
	services.CreateEngineService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("ENGINES")), redisClient)
	services.CreateAirlineService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("AIRLINES")), redisClient)
	services.CreateFlightService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("FLIGHTS")), redisClient)
	services.CreateReviewService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("REVIEWS")), redisClient)
	services.CreateRouteService(client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("ROUTES")), redisClient)

	// Routing
	// create a router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sessions.Sessions("x-airlines_api", store))
	authorized := router.Group("/")
	authorized.Use(AuthService.AuthMiddleware())
	{
		// users
		authorized.GET("/users", userController.GetUsers)
		authorized.GET("/users/airline_filter", userController.GetUserByAirline)
		authorized.GET("/users/:id/get_airlines", userController.GetUserAirlinesData)
		authorized.PUT("/users/:id", userController.UpdateUser)
		authorized.PUT("/users/:id/update_airlines", userController.UpdateUserAirlines)
		authorized.DELETE("/users/:id", userController.DeleteUser)

		// aircraft
		authorized.GET("/aircraft", aircraftController.GetAircraft)
		authorized.GET("/aircraft/:id", aircraftController.GetAircraftById)
		authorized.GET("/aircraft/aircraft_filter", aircraftController.GetAircraftByType)
		authorized.POST("/aircraft", aircraftController.CreateAircraft)
		authorized.DELETE("/aircraft/:id", aircraftController.DeleteAircraft)
		authorized.PUT("/aircraft/:id/update_airframe", aircraftController.UpdateAirframe)
		authorized.PUT("/aircraft/:id/update_exterior", aircraftController.UpdateExterior)
		authorized.PUT("/aircraft/:id/update_interior", aircraftController.UpdateInterior)
		authorized.PUT("/aircraft/:id/update_engines", aircraftController.UpdateEngines)
		authorized.PUT("/aircraft/:id/update_cockpit", aircraftController.UpdateCockpit)
		authorized.PUT("/aircraft/:id/update_general", aircraftController.UpdateGeneral)
		authorized.PUT("/aircraft/:id/update_performance", aircraftController.UpdatePerformance)
		authorized.PUT("/aircraft/:id/update_tags", aircraftController.UpdateTags)
		authorized.PUT("/aircraft/:id/update_owner", aircraftController.UpdateOwner)
		authorized.GET("/aircraft/:id/get_owner", aircraftController.GetOwnerData)
		authorized.GET("/aircraft/:id/get_engines", aircraftController.GetEngineData)
		authorized.GET("/aircraft/:id/get_airline", aircraftController.GetAirlineData)

		// engines
		authorized.GET("/engines", engineController.GetEngines)
		authorized.GET("/engines/:id", engineController.GetEngineById)
		authorized.POST("/engines", engineController.CreateEngine)
		authorized.PUT("/engines/:id", engineController.UpdateEngine)
		authorized.DELETE("/engines/:id", engineController.DeleteEngine)

		// airlines
		authorized.GET("/airlines", airlineController.GetAirlines)
		authorized.POST("/airlines", airlineController.CreateAirline)
		authorized.DELETE("/airlines/:id", airlineController.DeleteAirline)

		authorized.PUT("/airlines/:id/update_general", airlineController.UpdateAirlineGeneral)
		authorized.PUT("/airlines/:id/update_review", airlineController.UpdateReviews)
		authorized.PUT("/airlines/:id/update_routes", airlineController.UpdateRoutes)
		authorized.PUT("/airlines/:id/update_fleet", airlineController.UpdateFleet)
		authorized.PUT("/airlines/:id/update_owner", airlineController.UpdateAirlineOwner)
		authorized.GET("/airlines/:id/get_fleet", airlineController.GetFleetData)
		authorized.GET("/airlines/:id/get_owner", airlineController.GetAirlineOwnerData)
	}

	// handlers
	// sign in/ sign out
	router.POST("/signup", AuthService.SignUp)
	router.POST("/signin", AuthService.SignIn)
	router.POST("/signout", AuthService.SignOut)
	router.POST("/refresh", AuthService.Refresh)
	router.GET("/", Welcome)

	// listen and serve
	router.Run(port)
}
