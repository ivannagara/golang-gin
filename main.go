package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ivannagara/golang-gin/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// var is like 'late' in dart, bcs it will be initialized later on call
// until now, i think that var is usually used for global variable, and will
// be initialized in each class as different individual instances
var ctx context.Context
var err error
var client *mongo.Client
var recipesHandler *handlers.RecipesHandler

func init() {
	// recipes = make([]Recipe, 0)
	// read the recipes.json file
	// ------------------------------------------------|
	// file, _ := os.ReadFile("recipes.json")          |
	//												   |
	// _ = json.Unmarshal([]byte(file), &recipes)      |
	// ------------------------------------------------|
	// the background context is an empty context that will be used
	// and does not have any deadline and is never cancelled.
	// Background context is usually used for initialization,  tests, main function,
	// and as the top-level context for upcoming requests.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	// The [Ping] function is used to check if the connection to the database is valid
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	// var listOfRecipes []interface{}
	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }
	// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print("Inserted recipes:", len(insertManyResult.InsertedIDs))
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.Run()
}

func Add(a int, b int) int {
	return a + b
}
