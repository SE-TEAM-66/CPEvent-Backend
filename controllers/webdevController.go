package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateWebDevSkill(c *gin.Context) {
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
		WebDev string
	}

	if c.BindJSON(&tecSkillsBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Create a new skill and associate it with the profile
	skill := models.Skill{
		ProfileID: uint(profileID),
		Type:      "WebDev",
	}

	// Save the skill to the database
	resultSkill := initializers.DB.Create(&skill)
	if resultSkill.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
		return
	}

	// Create associated entry for WebDev
	webDev := models.WebDev{
		SkillID: uint(skill.ID),
		WebDev:       tecSkillsBody.WebDev,
	}

	// Save WebDev to the database
	resultWebDev := initializers.DB.Create(&webDev)
	if resultWebDev.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create WebDev"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tecSkills": webDev})
}

func EditWebDevSkill(c *gin.Context) {
    // Get profile ID, skill ID, and WebDevSkill ID from the request parameters
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

    // Check if the WebDevSkill exists and belongs to the correct profile's skill
    var webDevSkill models.WebDev
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&webDevSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "WebDev skill not found or does not belong to the profile"})
        return
    }

    // Get WebDev skill details from the request body
    var webDevSkillBody struct {
        WebDev string `json:"webDev" binding:"required"`
    }

    if c.BindJSON(&webDevSkillBody) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Update the WebDev skill details
    webDevSkill.WebDev = webDevSkillBody.WebDev

    // Save the updated WebDev skill to the database
    if err := initializers.DB.Save(&webDevSkill).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update WebDev skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"webDevSkill": webDevSkill})
}

func DeleteWebDevSkill(c *gin.Context) {
    // Get profile ID and web development skill ID from the request parameters
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

    // Check if the web development skill exists and belongs to the correct profile's skill
    var webDevSkill models.WebDev
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&webDevSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Web development skill not found or does not belong to the profile"})
        return
    }

    // Delete the web development skill from the database
    if err := initializers.DB.Delete(&webDevSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete web development skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Web development skill deleted successfully"})
}

func GetAllWebDevSkills(c *gin.Context) {
	// Get profile ID from the request parameters
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

	// Query web development skills associated with the profile
	var webDevSkills []models.WebDev
	if err := initializers.DB.Where("skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", profileID).Find(&webDevSkills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve web development skills"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"webDevSkills": webDevSkills})
}
