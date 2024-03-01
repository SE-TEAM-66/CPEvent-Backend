package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateDataAnaSkill(c *gin.Context) {
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
			DataAna       string
		}

		if c.BindJSON(&tecSkillsBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		// Create a new skill and associate it with the profile
		skill := models.Skill{
			ProfileID: uint(profileID),
			Type : "DataAnalysis",
		}

		// Save the skill to the database
		resultSkill := initializers.DB.Create(&skill)
		if resultSkill.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
			return
		}

		// Create associated  for DataAna
		dataAna := models.DataAna{
			SkillID: uint(skill.ID),
			DataAna: tecSkillsBody.DataAna,
		}

		// Save DataAna to the database
		resultDataAna := initializers.DB.Create(&dataAna)
		if resultDataAna.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create DataAna"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tecSkills": dataAna})

}
