package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Car struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var Cars []Car

func main() {
	Cars = []Car{
		{ID: "101", Title: "Mercedes", Color: "Blue"},
		{ID: "202", Title: "Audi", Color: "Black"},
	}

	router := gin.Default()
	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCarByID)
	router.POST("/cars", createCar)
	router.DELETE("/cars/:id", deleteCar)

	router.Run("localhost:8080")
}

// createCar handles the creation of a new car
func createCar(c *gin.Context) {
	var newCar Car
	// Bind the JSON body to newCar struct and handle any errors
	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to bind JSON: " + err.Error()})
		return
	}

	// Append the new car to the Cars slice
	Cars = append(Cars, newCar)
	c.IndentedJSON(http.StatusCreated, newCar)
}

// getCars handles the retrieval of all cars
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Cars)
}

// getCarByID handles the retrieval of a single car by its ID
func getCarByID(c *gin.Context) {
	id := c.Param("id")
	// Iterate over the Cars slice to find the car with the given ID
	for _, car := range Cars {
		if car.ID == id {
			c.IndentedJSON(http.StatusOK, car)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car with ID " + id + " not found"})
}

// deleteCar handles the deletion of a car by its ID
func deleteCar(c *gin.Context) {
	id := c.Param("id")
	// Iterate over the Cars slice to find the car with the given ID
	for index, car := range Cars {
		if car.ID == id {
			// Remove the car from the slice
			Cars = append(Cars[:index], Cars[index+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Car with ID " + id + " deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car with ID " + id + " not found"})
}
