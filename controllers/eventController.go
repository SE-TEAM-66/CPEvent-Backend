package controllers

import (
	"net/http"

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
	result := initializers.DB.Preload("Groups").Find(&events)

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