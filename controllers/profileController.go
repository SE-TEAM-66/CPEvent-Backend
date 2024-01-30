package controllers

import (
	"net/http"

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
	//get id
	id := c.Param(":id")
	
	//get the profile
	var profile []models.Profile
	initializers.DB.First(&profile,id)


	c.JSON(200, gin.H{"profile":profile,})
}

func ProfileUpdate(c *gin.Context){
	//get id 
	id := c.Param(":id")

	//get data off req body
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
	// Get id
	id := c.Param("id")

	// Start a transaction
	tx := initializers.DB.Begin()

	// Defer the rollback in case of an error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete profile
	if err := tx.Delete(&models.Profile{}, id).Error; err != nil {
		// Rollback the transaction on error
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to delete profile"})
		return
	}

	// Commit the transaction
	tx.Commit()

	c.Status(200)
}
