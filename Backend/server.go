package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const port = ":3000"

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	connectionString := os.Getenv("CONNECTION_STRING")

	// mongo connection
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

	// handler of web requests from the frontend. Paras (request address, fnc to handle a request)
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(res, "City Driver Backend")
		if err != nil {
			fmt.Println(err)
		}
	})

	// start a web server listening a port 3000 with no handler. It ha to be called after a handler of requests
	fmt.Println(fmt.Sprintf("A server is listening a port %s", port))
	_ = http.ListenAndServe(port, nil) // _ means an empty holder for a var because http.ListenAndServe() returns error which we don't care now

}
