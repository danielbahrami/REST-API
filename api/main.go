package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
)

type user struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var users []user

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func createUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	if len(newUser.ID) == 0 || len(newUser.FirstName) == 0 || len(newUser.LastName) == 0 || len(newUser.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user details required"})
		return
	}
	for _, user := range users {
		if newUser.ID == user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
	}
	users = append(users, newUser)
	c.JSON(http.StatusCreated, gin.H{"success": "user created"})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser user
	if err := c.BindJSON(&updatedUser); err != nil {
		return
	}
	for i := range users {
		if users[i].ID == id {
			if len(updatedUser.FirstName) == 0 || len(updatedUser.LastName) == 0 || len(updatedUser.Email) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user details required"})
				return
			}
			users[i].FirstName = updatedUser.FirstName
			users[i].LastName = updatedUser.LastName
			users[i].Email = updatedUser.Email
			c.JSON(http.StatusOK, gin.H{"success": "user updated"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"success": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
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
