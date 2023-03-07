package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
}

var mongoUri = os.Getenv("MONGO_URI")

var client *mongo.Client
var db *mongo.Database

func init() {
	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	db = client.Database(os.Getenv("MONGO_DB"))
}

func getUsers(c *gin.Context) {
	usersCollection := db.Collection("users")

	cur, err := usersCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []User
	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users = append(users, elem)
	}

	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")

	// Convert the id to a primitive.ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersCollection := db.Collection("users")
	result := usersCollection.FindOne(context.TODO(), bson.M{"_id": objId})

	var user User
	if err := result.Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		return
	}

	if len(user.FirstName) == 0 || len(user.LastName) == 0 || len(user.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User details required"})
		return
	}

	usersCollection := db.Collection("users")
	result, err := usersCollection.InsertOne(context.TODO(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		return
	}

	if len(updatedUser.FirstName) == 0 || len(updatedUser.LastName) == 0 || len(updatedUser.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User details required"})
		return
	}

	usersCollection := db.Collection("users")
	result, err := usersCollection.UpdateOne(context.TODO(), bson.M{"_id": objId}, bson.D{
		{"$set", updatedUser},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount >= 1 {
		c.JSON(http.StatusOK, gin.H{"message": "User updated"})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersCollection := db.Collection("users")
	result, err := usersCollection.DeleteOne(context.TODO(), bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount >= 1 {
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserById)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
