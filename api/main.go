package main

import (
  "log"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"

  "github.com/gin-gonic/gin"
)

const DbName = "data.db"

type Impl struct {
	DB *gorm.DB
}

func (i *Impl) InitDB() {
	var err error
  i.DB, err = gorm.Open("sqlite3", DbName)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Status{}, &Participation{}, &Event{})
}

func (i *Impl) Close() {
  i.DB.Close()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	i := Impl{}
	i.InitDB()
  i.InitSchema()
  //defer i.Close()

	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		// v1.POST("/status", i.PostStatus)
		// v1.GET("/status", i.GetStatus)
		// v1.GET("/status/:id", i.GetStatus)
		// v1.PUT("/status/:id", i.UpdateStatus)
		// v1.DELETE("/status/:id", i.DeleteStatus)

		// v1.POST("/participations", i.PostParticipation)
		// v1.GET("/participations", i.GetParticipations)
		// v1.GET("/participations/:id", i.GetParticipation)
		// v1.PUT("/participations/:id", i.UpdateParticipation)
		// v1.DELETE("/participations/:id", i.DeleteParticipation)

		v1.POST("/events", i.PostEvent)
		v1.GET("/events", i.GetEvents)
		v1.GET("/events/:id", i.GetEvent)
		v1.PUT("/events/:id", i.UpdateEvent)
		v1.DELETE("/events/:id", i.DeleteEvent)
	}

  r.Run(":8080")
}