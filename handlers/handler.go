package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ivannagara/golang-gin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// (#) Why is the gin.Context required the be given a pointer?
// [*] Because, the context needs to change value and by using the pointer,
// the function can change its own value.
// Here, we can see that the 'NewRecipeHandler' function needs to be able to
// set the status code value, and the response body of the HTTP response,
// which are a part of the context object.
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Remove data from Redis")
	handler.redisClient.Del("recipes")
	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)
		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe) // should have &
			recipes = append(recipes, recipe)
		}
		data, _ := json.Marshal(recipes)
		handler.redisClient.Set("recipes", string(data), 0)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to REDIS")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}
}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectID}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: recipe.Name},
		{Key: "instructions", Value: recipe.Instructions},
		{Key: "ingredients", Value: recipe.Ingredients},
		{Key: "tags", Value: recipe.Tags},
	}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	log.Println("Remove data from Redis")
	handler.redisClient.Del("recipes")
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// func DeleteRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	index := -1
// 	for i := 0; i < len(recipes); i++ {
// 		if recipes[i].ID == id {
// 			index = i
// 		}
// 	}
// 	if index == -1 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Recipe not found"})
// 		return
// 	}
// 	recipes = append(recipes[:index], recipes[index+1:]...)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Recipe has been deleted"})
// }
