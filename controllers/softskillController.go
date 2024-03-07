package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateSoftSkill(c *gin.Context) {
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

		// Get soft skill details from the request body
		var softSkillBody struct {
			Title string
		}

		if c.BindJSON(&softSkillBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		// Create a new skill and associate it with the profile
		skill := models.Skill{
			ProfileID: uint(profileID),
			Type:"Soft_skill",
		}

		// Save the skill to the database
		resultSkill := initializers.DB.Create(&skill)
		if resultSkill.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
			return
		}

		// Create a new soft skill and associate it with the skill
		softSkill := models.Soft_skill{
			SkillID: uint(skill.ID), // Assuming Skill has an ID field
			Title:   softSkillBody.Title,

		}

		// Save the soft skill to the database
		resultSoftSkill := initializers.DB.Create(&softSkill)
		if resultSoftSkill.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create soft skill"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"softSkill": softSkill})
}

func EditSoftSkill(c *gin.Context) {
    // Get profile ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID
    var profile models.Profile
    if err := initializers.DB.Preload("Skill.Soft_skill").First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Check if there are any soft skills
    if len(profile.Skill) == 0 || len(profile.Skill[0].Soft_skill) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No soft skills found for the profile"})
        return
    }

    // Parse the request body to get the new title for the first soft skill
    var requestBody struct {
        Title string `json:"title"`
    }
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Update the title of the first soft skill
    firstSoftSkill := profile.Skill[0].Soft_skill[0]
    firstSoftSkill.Title = requestBody.Title
    if err := initializers.DB.Save(&firstSoftSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update soft skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Soft skill updated successfully"})
}



func DeleteSoftSkill(c *gin.Context) {
    // Get profile ID and soft skill ID from the request parameters
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

    // Check if the soft skill exists and belongs to the correct profile's skill
    var softSkill models.Soft_skill
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&softSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Soft skill not found or does not belong to the profile"})
        return
    }

    // Delete the soft skill from the database
    if err := initializers.DB.Delete(&softSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete soft skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Soft skill deleted successfully"})
}

func GetAllSoftSkills(c *gin.Context) {
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

    // Query soft skills associated with the profile
    var softSkills []models.Soft_skill
    if err := initializers.DB.Model(&models.Soft_skill{}).Select("title").Joins("JOIN skills ON soft_skills.skill_id = skills.id").Where("skills.profile_id = ?", profileID).Find(&softSkills).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve soft skills"})
        return
    }

    // Extract titles from soft skills
    var titles []string
    for _, softSkill := range softSkills {
        titles = append(titles, softSkill.Title)
    }

    c.JSON(http.StatusOK, gin.H{"softSkills": titles})
}
