package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateProgrammingSkill(c *gin.Context) {
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
		Programtype string
	}

	if c.BindJSON(&tecSkillsBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Create a new skill and associate it with the profile
	skill := models.Skill{
		ProfileID: uint(profileID),
		Type:      "TechnicalSkill",
	}

	// Save the skill to the database
	resultSkill := initializers.DB.Create(&skill)
	if resultSkill.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
		return
	}

	// Create a new technical skills entry and associate it with the skill
	tecSkills := models.Tec_skills{
		SkillID: uint(skill.ID),
		Type:    "Programming",
	}

	// Save the technical skills entry to the database
	resultTecSkills := initializers.DB.Create(&tecSkills)
	if resultTecSkills.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create technical skills"})
		return
	}

	// Create associated entry for Programming
	programming := models.Programming{
		Tec_skillsID: uint(tecSkills.ID),
		Programtype: tecSkillsBody.Programtype,
	}

	// Save Programming to the database
	resultProgramming := initializers.DB.Create(&programming)
	if resultProgramming.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Programming"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tecSkills": tecSkills})
}
