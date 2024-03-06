package controllers

import (
	"net/http"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func AddMember(c *gin.Context) {
	// Get data from request body
	var requestBody struct {
		GroupID uint   `json:"group_id" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the group
	var group models.Group
	if err := initializers.DB.First(&group, requestBody.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Find the user based on the provided email
	var user models.User
	if err := initializers.DB.Where("email = ?", requestBody.Email).Preload("Profile").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user is already a member of the group
	var existingMember models.Member
	if err := initializers.DB.Where("profile_id = ?", user.Profile.ID).Where("group_id = ?", group.ID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member of the group"})
		return
	}

	// Create a new member and associate it with the group
	newMember := models.Member{
		Role:      "", // You can set the default role for the member
		ProfileID: user.Profile.ID,
		GroupID:   group.ID,
	}

	if err := initializers.DB.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add member to the group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

func DeleteMember(c *gin.Context) {
	// Get data from request body
	var requestBody struct {
		GroupID   uint `json:"group_id" binding:"required"`
		ProfileID uint `json:"profile_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the group
	var group models.Group
	if err := initializers.DB.First(&group, requestBody.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Find the member to delete
	var memberToDelete models.Member
	if err := initializers.DB.
		Where("group_id = ?", group.ID).
		Where("profile_id = ?", requestBody.ProfileID).
		First(&memberToDelete).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// Delete the member
	if err := initializers.DB.Delete(&memberToDelete).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}
