package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var people = []person{
	{ID: "1", FirstName: "Emalee", LastName: "Creigan", Email: "ecreigan0@nature.com"},
	{ID: "2", FirstName: "Mariellen", LastName: "Peerless", Email: "mpeerless1@dot.gov"},
	{ID: "3", FirstName: "Gabie", LastName: "Brims", Email: "gbrims2@cnbc.com"},
	{ID: "4", FirstName: "Francesco", LastName: "Morman", Email: "fmorman3@rediff.com"},
	{ID: "5", FirstName: "Ario", LastName: "Denerley", Email: "adenerley4@state.gov"},
}

func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, people)
}

func personById(c *gin.Context) {
	id := c.Param("id")
	person, err := getPersonById(id)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, person)
}

func getPersonById(id string) (*person, error) {
	for i, person := range people {
		if person.ID == id {
			return &people[i], nil
		}
	}
	return nil, errors.New("person not found")
}

func addPerson(c *gin.Context) {

	var newPerson person

	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	people = append(people, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)

}

func main() {
	router := gin.Default()
	router.GET("/people", getPeople)
	router.GET("/people/:id", personById)
	router.POST("/people", addPerson)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
