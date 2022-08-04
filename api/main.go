package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type user struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var users = []user{
	{ID: "1", FirstName: "Emalee", LastName: "Creigan", Email: "ecreigan0@nature.com"},
	{ID: "2", FirstName: "Mariellen", LastName: "Peerless", Email: "mpeerless1@dot.gov"},
	{ID: "3", FirstName: "Gabie", LastName: "Brims", Email: "gbrims2@cnbc.com"},
	{ID: "4", FirstName: "Francesco", LastName: "Morman", Email: "fmorman3@rediff.com"},
	{ID: "5", FirstName: "Ario", LastName: "Denerley", Email: "adenerley4@state.gov"},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.IndentedJSON(http.StatusOK, user)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func createUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserById)
	router.POST("/users", createUser)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
