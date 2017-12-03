package main

import (
  "time"
  "net/http"
  "github.com/jinzhu/gorm"
  "github.com/gin-gonic/gin"
)

type Event struct {
  gorm.Model

  Name            string          `json:"name" binding:"required"`
  Date            time.Time       `json:"date" binding:"required"`

  Status          Status
  StatusID        int             `json:"status-id" binding:"required"`

  Participations  []Participation
}

// http --form POST :8080/api/v1/events name="My Event" description="My Description"
func (i *Impl) PostEvent(c *gin.Context) {
	var event Event
  if err := c.Bind(&event); err == nil {
    i.DB.Create(&event)
    c.JSON(http.StatusOK, gin.H{"success": event})
  } else {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }
}

// http :8080/api/v1/events
func (i *Impl) GetEvents(c *gin.Context) {
  var events []Event
  i.DB.Find(&events)
	c.JSON(http.StatusOK, events)
}

// http :8080/api/v1/events/1
func (i *Impl) GetEvent(c *gin.Context) {
  id := c.Param("id")
  var event Event
  if i.DB.First(&event, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	} else {
    c.JSON(http.StatusOK, event)
  }
}

// http PUT :8080/api/v1/events/1  name="My Event Updated" description="My Description Updated"
func (i *Impl) UpdateEvent(c *gin.Context) {
  var event Event
  var status Status
  if err := c.Bind(&event); err == nil {
    if i.DB.First(&event, event.ID).Error != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
    } else if i.DB.First(&status, event.StatusID).Error != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
    } else {
      i.DB.Save(&event)
      c.JSON(http.StatusOK, gin.H{"success": event})
    }
  } else {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }  
}

// http DELETE :8080/api/v1/events/1
func (i *Impl) DeleteEvent(c *gin.Context) {
  id := c.Param("id")
	var event Event  
  if i.DB.First(&event, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	} else {
    i.DB.Delete(&event)
    c.JSON(http.StatusOK, gin.H{"success": event})
  }
}