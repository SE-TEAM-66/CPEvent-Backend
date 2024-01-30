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
		ProfilePicture:"some url",
		Fname:"nica",
		Lname:"ragua",
		Faculty: "CPE",
		Bio:"i love chicken wings",
		Phone:"02345678",
		Email: "nicaragua111@gmail.com",
		Facebook: "some urllll",
		Line: "ID BLABLABLA"} 

	result := initializers.DB.Create(&profile)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, gin.H{"profile":profile,})
}

// func ProfileIndex(c *gin.Context){
// 	var posts []models.
// 	initializers.DB.Find(&posts)

// }