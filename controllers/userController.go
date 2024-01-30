package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func User(c *gin.Context) {
	//Get
	var body struct {
		Fname    string
		Lname    string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	user := models.User{
        Fname: body.Fname,
        Lname: body.Lname,
        Email: body.Email,
        Profile: models.Profile{
            Fname: body.Fname, 
            Lname: body.Lname,
        },
    }

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}
	
	// Fetch the associated profile from the database
	var profile models.Profile
	err := initializers.DB.Model(&user).Association("Profile").Find(&profile)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetProfileWithUser(c *gin.Context) {
    // Get profile ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Fetch the profile along with its associated user 
    var profile models.Profile
    if err := initializers.DB.Preload("User").Where("id = ?", profileID).First(&profile).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Return the profile and its user
    c.JSON(http.StatusOK, gin.H{"profile": profile})
}
