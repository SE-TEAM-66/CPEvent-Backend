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
		Type:      "Programming",
	}

	// Save the skill to the database
	resultSkill := initializers.DB.Create(&skill)
	if resultSkill.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
		return
	}

	// Create associated entry for Programming
	programming := models.Programming{
		SkillID: uint(skill.ID),
		Programtype: tecSkillsBody.Programtype,
	}

	// Save Programming to the database
	resultProgramming := initializers.DB.Create(&programming)
	if resultProgramming.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Programming"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tecSkills": programming})
}

func EditProgrammingSkill(c *gin.Context) {
    // Get profile ID, skill ID, and ProgrammingSkill ID from the request parameters
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

    // Check if the ProgrammingSkill exists and belongs to the correct profile's skill
    var programmingSkill models.Programming
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&programmingSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Programming skill not found or does not belong to the profile"})
        return
    }

    // Get Programming skill details from the request body
    var programmingSkillBody struct {
        Programtype string `json:"programtype" binding:"required"`
    }

    if c.BindJSON(&programmingSkillBody) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Update the Programming skill details
    programmingSkill.Programtype = programmingSkillBody.Programtype

    // Save the updated Programming skill to the database
    if err := initializers.DB.Save(&programmingSkill).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update Programming skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"programmingSkill": programmingSkill})
}

func DeleteProgrammingSkill(c *gin.Context) {
    // Get profile ID and programming skill ID from the request parameters
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

    // Check if the programming skill exists and belongs to the correct profile's skill
    var programmingSkill models.Programming
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&programmingSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Programming skill not found or does not belong to the profile"})
        return
    }

    // Delete the programming skill from the database
    if err := initializers.DB.Delete(&programmingSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete programming skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Programming skill deleted successfully"})
}
