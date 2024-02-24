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

func ProfileIndex(c *gin.Context){
	//get the profile
	var profiles []models.Profile
	initializers.DB.Find(&profiles)


	c.JSON(200, gin.H{"profiles":profiles,})
}

func ProfileShow(c *gin.Context){
	// Get id from URL parameter
	id := c.Param("id")
	
	// Get the profile from the database by ID
	var profile models.Profile
	if err := initializers.DB.First(&profile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Return the profile in the response
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}


func ProfileUpdate(c *gin.Context) {
	// Get id
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	if err := initializers.DB.First(&profile, id).Error; err != nil {
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
	id := c.Param("id")

	// Check if the profile exists
	var profile models.Profile
	if err := initializers.DB.First(&profile, id).Error; err != nil {
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

