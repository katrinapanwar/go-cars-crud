package main

import (
	"database/sql"
	f "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Car represents a car in the database
type Car struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var db *sql.DB

func main() {
	var err error
	// Connect to the MySQL database
	db, err = sql.Open("mysql", "root:admin@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		f.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}

	router := gin.Default()
	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCarByID)
	router.POST("/cars", createCar)
	router.PUT("/cars/:id", updateCar)
	router.DELETE("/cars/:id", deleteCar)

	router.Run("localhost:8080")
}

// getCars handles the retrieval of all cars
func getCars(c *gin.Context) {
	stmt := "SELECT * FROM cars"
	rows, err := db.Query(stmt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Title, &car.Color); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cars = append(cars, car)
	}
	c.JSON(http.StatusOK, cars)
}

// getCarByID handles the retrieval of a single car by its ID
func getCarByID(c *gin.Context) {
	id := c.Param("id")
	stmt := "SELECT * FROM cars WHERE id = ?"
	row := db.QueryRow(stmt, id)

	var car Car
	if err := row.Scan(&car.ID, &car.Title, &car.Color); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, car)
}

// createCar handles the creation of a new car
func createCar(c *gin.Context) {
	var newCar Car
	if err := c.BindJSON(&newCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind JSON: " + err.Error()})
		return
	}

	stmt := "INSERT INTO cars (title, color) VALUES (?, ?)"
	res, err := db.Exec(stmt, newCar.Title, newCar.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newCar.ID = int(id)
	c.JSON(http.StatusCreated, newCar)
}

// updateCar handles the updating of an existing car by its ID
func updateCar(c *gin.Context) {
	id := c.Param("id")
	var updatedCar Car
	if err := c.BindJSON(&updatedCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind JSON: " + err.Error()})
		return
	}

	stmt := "UPDATE cars SET title = ?, color = ? WHERE id = ?"
	res, err := db.Exec(stmt, updatedCar.Title, updatedCar.Color, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car successfully updated"})
}

// deleteCar handles the deletion of a car by its ID
func deleteCar(c *gin.Context) {
	id := c.Param("id")

	stmt := "DELETE FROM cars WHERE id = ?"
	res, err := db.Exec(stmt, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car successfully deleted"})
}

/*package main

import (
	"database/sql"
	f "fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open may just validate its arguments without creating a connection to the database.
	// To verify that the data source name is valid, call Ping.
	// func Open(driverName, dataSourceName string) (*DB, error)
	// dataSourceName: username:password@protocol(address)/dbname?param=value
	//db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")

	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/testdb")

	if err != nil {
		f.Println("error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		f.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}

	// func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	insert, err := db.Query("INSERT INTO `testdb`.`students` (`StudentID`, `FirstName`, `LastName`) VALUES ('3', 'Carl', 'Jones');")
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
	f.Println("Successful Connection to Database!")
}

*/
