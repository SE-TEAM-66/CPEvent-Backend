package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func ProfileCreate(c *gin.Context) {
	//Get
	var body struct {
		ProfilePicture string
		Fname string
		Lname string
		Faculty string
		Bio string
		Phone string
		Email string
		Facebook string
		Line string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	//Create profile
	profile := models.Profile{
		ProfilePicture:body.ProfilePicture,
		Fname:body.Fname,
		Lname:body.Lname,
		Faculty: body.Faculty,
		Bio:body.Bio,
		Phone:body.Phone,
		Email:body.Email,
		Facebook:body.Facebook,
		Line:body.Line} 

	result := initializers.DB.Create(&profile)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, gin.H{"profile":profile,})
}

func ProfileIndex(c *gin.Context) {
    var profiles []models.Profile

    // Find all profiles including their associated experiences
    if err := initializers.DB.Preload("Exp").Find(&profiles).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

func ProfileShow(c *gin.Context) {
    // Get profile ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID including its associated experiences
    var profile models.Profile
    if err := initializers.DB.Preload("Exp").First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"profile": profile})
}



func ProfileUpdate(c *gin.Context) {
	// Get id
	profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	// Get data from request body
	var body struct {
		ProfilePicture string
		Fname          string
		Lname          string
		Faculty        string
		Bio            string
		Phone          string
		Email          string
		Facebook       string
		Line           string
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Find the profile to be updated
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Update the profile
	updatedProfile := models.Profile{
		ProfilePicture: body.ProfilePicture,
		Fname:          body.Fname,
		Lname:          body.Lname,
		Faculty:        body.Faculty,
		Bio:            body.Bio,
		Phone:          body.Phone,
		Email:          body.Email,
		Facebook:       body.Facebook,
		Line:           body.Line,
	}
	if err := initializers.DB.Model(&profile).Updates(updatedProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": updatedProfile})
}

func ProfileDelete(c *gin.Context) {
	// Get ID from URL parameter
	profileID := c.Param("profileID")

	// Check if the profile exists
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Start a transaction
	tx := initializers.DB.Begin()

	// Defer the rollback in case of an error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete profile
	if err := tx.Delete(&profile).Error; err != nil {
		// Rollback the transaction on error
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	// Commit the transaction
	tx.Commit()

	c.Status(http.StatusOK)
}


func ProfileUpdate(c *gin.Context){
	//get id 
	id := c.Param(":id")

	//get data off req body
	var body struct {
		ProfilePicture string
		Fname          string
		Lname          string
		Faculty        string
		Bio            string
		Phone          string
		Email          string
		Facebook       string
		Line           string
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	//find post which want to be updated
	var profile []models.Profile
	initializers.DB.First(&profile,id)

	//update
	initializers.DB.Model(&profile).Updates(models.Profile{
		ProfilePicture:body.ProfilePicture,
		Fname:body.Fname,
		Lname:body.Lname,
		Faculty: body.Faculty,
		Bio:body.Bio,
		Phone:body.Phone,
		Email:body.Email,
		Facebook:body.Facebook,
		Line:body.Line,
	})
	
	c.JSON(200, gin.H{"profile":profile,})
}

func ProfileDelete(c *gin.Context) {
	// Get ID from URL parameter
	profileID := c.Param("profileID")

	// Check if the profile exists
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Start a transaction
	tx := initializers.DB.Begin()

	// Defer the rollback in case of an error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete profile
	if err := tx.Delete(&profile).Error; err != nil {
		// Rollback the transaction on error
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	// Commit the transaction
	tx.Commit()

	c.Status(http.StatusOK)
}

func ProfileImage(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in the context"})
		return
	}

	// Type assertion
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type in the context"})
		return
	}

	// Fetch the profile using the user's ID
	var profile models.Profile
	result := initializers.DB.Where("user_id = ?", userModel.ID).Find(&profile)

	if result.Error != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Process the retrieved profile
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
