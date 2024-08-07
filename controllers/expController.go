package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func GetAllExperiences(c *gin.Context) {
	// Query all experiences
	var experiences []models.Exp
	if err := initializers.DB.Find(&experiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve experiences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"experiences": experiences})
}


func GetExperience(c *gin.Context) {
    // Get profile ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID
    var profile models.Profile
    if err := initializers.DB.First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Query experiences associated with the profile ID
    var experiences []models.Exp
    if err := initializers.DB.Where("profile_id = ?", profile.ID).Find(&experiences).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve experiences"})
        return
    }

    // Create a custom response structure to include only Title and Description
    var response []gin.H
    for _, exp := range experiences {
        response = append(response, gin.H{
            "Title":       exp.Title,
            "Description": exp.Description,
        })
    }

    c.JSON(http.StatusOK, response)
}


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

	// associate the created group to the owner id
	initializers.DB.Model(&experience).Association("Profiles").Append(&profile)

	c.JSON(http.StatusOK, gin.H{"experience": experience})
}

func EditExperience(c *gin.Context) {
    // Get profile ID and experience ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID
    var profile models.Profile
    if err := initializers.DB.Preload("Exp").First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Check if there are any experiences
    if len(profile.Exp) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No experiences found for the profile"})
        return
    }

    // Parse the request body to get the new title and description for the experience
    var requestBody struct {
        Description string `json:"description"`
    }
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Update the title and description of the first experience
    firstExperience := profile.Exp[0]
    firstExperience.Description = requestBody.Description

    // Save the updated experience to the database
    if err := initializers.DB.Save(&firstExperience).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update experience"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Experience updated successfully"})
}


func DeleteExperience(c *gin.Context) {
    // Get profile ID and experience ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID
    var profile models.Profile
    if err := initializers.DB.First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    expID, err := strconv.ParseUint(c.Param("expID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ExpID"})
        return
    }

    // Check if the experience exists and belongs to the correct profile
    var experience models.Exp
    if err := initializers.DB.Where("profile_id = ? AND id = ?", profileID, expID).First(&experience).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found or does not belong to the profile"})
        return
    }

    // Delete the experience
    if err := initializers.DB.Delete(&experience).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete experience"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Experience deleted successfully"})
}
