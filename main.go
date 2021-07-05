// Recipes API
//
// This is a sample API covering fairly simple endpoints in a
// recipes domain.
//
// Schemes: http
// Host: localhost:8080
// BASEPATH: /
// Version: 1.0.0
// Contact: Ed Leonard <edward.leonard@gmail.com>
//
// consumes:
// -application/json
//
// Produces:
// -application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	handlers "go-recipes-api/handlers"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

//var recipes []Recipe

//mongoDB
//var ctx context.Context
//var err error

//var client *mongo.Client

func init() {
	// mongoDB client connection
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)

	log.Println("Successfully Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	// routes
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.Run()
}
