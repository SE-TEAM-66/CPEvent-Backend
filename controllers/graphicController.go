package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateGraphicDesignSkill(c *gin.Context) {
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

	// Get technical skills details from the request body
	var tecSkillsBody struct {
		GraphicDesign string
	}

	if c.BindJSON(&tecSkillsBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Create a new skill and associate it with the profile
	skill := models.Skill{
		ProfileID: uint(profileID),
		Type: "GraphicDesign",
	}

	// Save the skill to the database
	resultSkill := initializers.DB.Create(&skill)
	if resultSkill.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
		return
	}

	// Create associated entry for GraphicDesign
	graphicDesign := models.GraphicDesign{
		SkillID:  uint(skill.ID),
		GraphicDesign: tecSkillsBody.GraphicDesign,
	}

	// Save GraphicDesign to the database
	resultGraphicDesign := initializers.DB.Create(&graphicDesign)
	if resultGraphicDesign.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create GraphicDesign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tecSkills": graphicDesign})
}
