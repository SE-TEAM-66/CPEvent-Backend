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

func EditDataAnaSkill(c *gin.Context) {
    // Get profile ID, skill ID, and DataAnaSkill ID from the request parameters
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

    skillID, err := strconv.ParseUint(c.Param("skillID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SkillID"})
        return
    }

    // Check if the DataAnaSkill exists and belongs to the correct profile's skill
    var dataAnaSkill models.DataAna
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&dataAnaSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "DataAna skill not found or does not belong to the profile"})
        return
    }

    // Get DataAna skill details from the request body
    var dataAnaSkillBody struct {
        DataAna string `json:"dataAna" binding:"required"`
    }

    if c.BindJSON(&dataAnaSkillBody) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Update the DataAna skill details
    dataAnaSkill.DataAna = dataAnaSkillBody.DataAna

    // Save the updated DataAna skill to the database
    if err := initializers.DB.Save(&dataAnaSkill).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update DataAna skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"dataAnaSkill": dataAnaSkill})
}