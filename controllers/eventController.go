package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func EventCreate(c *gin.Context) {
	// Get data from req body
	var body struct {
		PicUrl string `json:"pic_url" binding:"required"`
		Etitle string `json:"e_title" binding:"required"`
		Edesc  string `json:"e_desc" binding:"required"`
		Edate  string `json:"e_date" binding:"required"`
		Etime  string `json:"e_time" binding:"required"`
	}

	// Bind and validate
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Create Event
	event := models.Event{
		PicUrl: body.PicUrl,
		Etitle: body.Etitle,
		Edesc: body.Edesc,
		Edate: body.Edate,
		Etime: body.Etime,
	}

	if err := initializers.DB.Create(&event); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create event",
		})
		return
	}

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

}

func GetAllEvents(c * gin.Context){
	//Get all events
	var events []models.Event
	result := initializers.DB.Preload("Groups").Preload("Groups.ReqPositions").Find(&events)

	//Return on error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to all events",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": events,
	})
}
func GetSingleEvent(c *gin.Context){
    // Get eid from param
    eidStr := c.Param("eid")
    eid, err := strconv.Atoi(eidStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
        return
    }

    // Find Event, pre-populate fields to optimize query
    var event models.Event
    result := initializers.DB.Preload("Groups").First(&event, eid)

    // Handle errors gracefully
    if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target event not found",
		})
		return	
    }

    // Success response, optionally format data or remove sensitive information
    c.JSON(http.StatusOK, gin.H{
		"message": event,
	})
}