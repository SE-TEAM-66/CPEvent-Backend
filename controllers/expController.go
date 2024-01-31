package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateExperience(c *gin.Context) {
	// Get profile ID from the request parameters
	profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
		return
	}

	// Check if the profile exists
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Get experience details from the request body
	var experienceBody struct {
		Title       string 
		Description string 
	}

	if c.BindJSON(&experienceBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Create a new experience and associate it with the profile
	experience := models.Exp{
		ProfileID:   uint(profileID),
		Title:       experienceBody.Title,
		Description: experienceBody.Description,
	}

	// Save the experience to the database
	result := initializers.DB.Create(&experience)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create experience"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"experience": experience})
}
