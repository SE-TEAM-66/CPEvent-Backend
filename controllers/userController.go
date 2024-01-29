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
		Phone   string
		Bio 	string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	user := models.User{Faculty:"Department of Computer Engineering",Phone:body.Phone, Bio:body.Bio} 

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, gin.H{"posts":user,})
}

// func ProfileIndex(c *gin.Context){
// 	var posts []models.
// 	initializers.DB.Find(&posts)

// }