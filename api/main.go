package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Status struct {
	Id     int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Name   string `gorm:"not null" form:"name" json:"name"`
	Color  string `gorm:"not null" form:"color" json:"color"`
}

type Events struct {
	Id           int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Name         string `gorm:"not null" form:"name" json:"name"`
	Description  string `gorm:"not null" form:"description" json:"description"`
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Events{}) {
		db.CreateTable(&Events{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Events{})
	}

	return db
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.POST("/events", PostEvent)
		v1.GET("/events", GetEvents)
		v1.GET("/events/:id", GetEvent)
		v1.PUT("/events/:id", UpdateEvent)
		v1.DELETE("/events/:id", DeleteEvent)
		v1.OPTIONS("/events", OptionsEvent)      // POST
		v1.OPTIONS("/events/:id", OptionsEvent)  // PUT, DELETE
	}

	r.Run(":8080")
}

func PostEvent(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var event Events
	c.Bind(&event)

	if event.Name != "" && event.Description != "" {
		// INSERT event
		db.Create(&event)
		// Display error
		c.JSON(201, gin.H{"success": event})
	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// http --form POST :8080/api/v1/events name="My Event" description="My Description"
}

func GetEvents(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var events []Events
	// SELECT events
	db.Find(&events)

	// Display JSON result
	c.JSON(200, events)

	// http :8080/api/v1/events
}

func GetEvent(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var event Events
	// SELECT event
	db.First(&event, id)

	if event.Id != 0 {
		// Display JSON result
		c.JSON(200, event)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Event not found"})
	}

	// http :8080/api/v1/events/1
}

func UpdateEvent(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id user
	id := c.Params.ByName("id")
	var event Events
	// SELECT * FROM events WHERE id = 1;
	db.First(&event, id)

	if event.Name != "" && event.Description != "" {

		if event.Id != 0 {
			var newEvent Events
			c.Bind(&newEvent)

			result := Events{
				Id:           event.Id,
				Name:         newEvent.Name,
				Description:  newEvent.Description,
			}

			// UPDATE event
			db.Save(&result)
			// Display modified data in JSON message "success"
			c.JSON(200, gin.H{"success": result})
		} else {
			// Display JSON error
			c.JSON(404, gin.H{"error": "Event not found"})
		}

	} else {
		// Display JSON error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// http PUT :8080/api/v1/events/1  name="My Event Updated" description="My Description Updated"
}

func DeleteEvent(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id event
	id := c.Params.ByName("id")
	var event Events
	// SELECT event
	db.First(&event, id)

	if event.Id != 0 {
		// DELETE event
		db.Delete(&event)
		// Display JSON result
		c.JSON(200, gin.H{"success": "Event #" + id + " deleted"})
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Event not found"})
	}

	// http DELETE :8080/api/v1/events/1
}

func OptionsEvent(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}