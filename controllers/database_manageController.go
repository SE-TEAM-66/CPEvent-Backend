package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateDBManageSkill(c *gin.Context) {
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
		DBManage string 
	}

	if c.BindJSON(&tecSkillsBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Create a new skill and associate it with the profile
	skill := models.Skill{
		ProfileID: uint(profileID),
		Type:      "DBManagement",
	}

	// Save the skill to the database
	resultSkill := initializers.DB.Create(&skill)
	if resultSkill.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
		return
	}


	// Create associated entry for DBManage
	dbManage := models.DBmanage{
		SkillID: uint(skill.ID),
		DBmanage: tecSkillsBody.DBManage,
	}

	// Save DBManage to the database
	resultDBManage := initializers.DB.Create(&dbManage)
	if resultDBManage.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create DBManage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tecSkills": dbManage})
}

func EditDBManageSkill(c *gin.Context) {
    // Get profile ID, skill ID, and DBManageSkill ID from the request parameters
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

    // Check if the DBManageSkill exists and belongs to the correct profile's skill
    var dbManageSkill models.DBmanage
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&dbManageSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "DBManage skill not found or does not belong to the profile"})
        return
    }

    // Get DBManage skill details from the request body
    var dbManageSkillBody struct {
        DBManage string `json:"dbManage" binding:"required"`
    }

    if c.BindJSON(&dbManageSkillBody) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Update the DBManage skill details
    dbManageSkill.DBmanage = dbManageSkillBody.DBManage

    // Save the updated DBManage skill to the database
    if err := initializers.DB.Save(&dbManageSkill).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update DBManage skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"dbManageSkill": dbManageSkill})
}

func DeleteDBManageSkill(c *gin.Context) {
    // Get profile ID and database management skill ID from the request parameters
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

    // Check if the database management skill exists and belongs to the correct profile's skill
    var dbManageSkill models.DBmanage
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&dbManageSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Database management skill not found or does not belong to the profile"})
        return
    }

    // Delete the database management skill from the database
    if err := initializers.DB.Delete(&dbManageSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete database management skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Database management skill deleted successfully"})
}