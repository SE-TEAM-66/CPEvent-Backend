package controllers

import (
	"fmt"
	"net/http"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func NotifyCreate(c *gin.Context) {
	//Get data from req body
	var body struct {
		Rec_id  uint
		Sender  string
		Message string //accept or not
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	notify := models.Notify{
		Rec_id:  body.Rec_id,
		Sender:  body.Sender,
		Message: body.Message,
		IsRead:  false,
	}

	if err := initializers.DB.Create(&notify); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create group",
		})
		return
	}

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": notify,
	})

}

func NotifyGet(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in the context!"})
		return
	}

	// Type assertion
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type in the context"})
		return
	}

	// Fetch the profile using the user's ID
	var notification []models.Notify
	result := initializers.DB.Where("rec_id = ?", userModel.Profile.ID).Order("id DESC").Find(&notification)
	fmt.Println(userModel)
	if result.Error != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, notification)
}
