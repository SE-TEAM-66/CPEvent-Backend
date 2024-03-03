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

func EditGraphicDesignSkill(c *gin.Context) {
    // Get profile ID, skill ID, and GraphicDesignSkill ID from the request parameters
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

    // Check if the GraphicDesignSkill exists and belongs to the correct profile's skill
    var graphicDesignSkill models.GraphicDesign
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&graphicDesignSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "GraphicDesign skill not found or does not belong to the profile"})
        return
    }

    // Get GraphicDesign skill details from the request body
    var graphicDesignSkillBody struct {
        GraphicDesign string `json:"graphicDesign" binding:"required"`
    }

    if c.BindJSON(&graphicDesignSkillBody) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Update the GraphicDesign skill details
    graphicDesignSkill.GraphicDesign = graphicDesignSkillBody.GraphicDesign

    // Save the updated GraphicDesign skill to the database
    if err := initializers.DB.Save(&graphicDesignSkill).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update GraphicDesign skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"graphicDesignSkill": graphicDesignSkill})
}

func DeleteGraphicDesignSkill(c *gin.Context) {
    // Get profile ID and graphic design skill ID from the request parameters
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

    // Check if the graphic design skill exists and belongs to the correct profile's skill
    var graphicDesignSkill models.GraphicDesign
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&graphicDesignSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Graphic design skill not found or does not belong to the profile"})
        return
    }

    // Delete the graphic design skill from the database
    if err := initializers.DB.Delete(&graphicDesignSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete graphic design skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Graphic design skill deleted successfully"})
}

func GetAllGraphicDesignSkills(c *gin.Context) {
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

	// Query graphic design skills associated with the profile
	var graphicDesignSkills []models.GraphicDesign
	if err := initializers.DB.Where("skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", profileID).Find(&graphicDesignSkills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve graphic design skills"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"graphicDesignSkills": graphicDesignSkills})
}