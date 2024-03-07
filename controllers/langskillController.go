package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func CreateLangSkill(c *gin.Context) {
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

// Get language skill details from the request body
var langSkillBody struct {
    Title string
}

if c.BindJSON(&langSkillBody) != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
    return
}

// Create a new skill and associate it with the profile
skill := models.Skill{
    ProfileID: uint(profileID),
	Type:"Lang_skill",
}

// Save the skill to the database
resultSkill := initializers.DB.Create(&skill)
if resultSkill.Error != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create skill"})
    return
}

// Create a new language skill and associate it with the skill
langSkill := models.Lang_skill{
    SkillID: uint(skill.ID), // Assuming Skill has an ID field
    Title:   langSkillBody.Title,
}

// Save the language skill to the database
resultLangSkill := initializers.DB.Create(&langSkill)
if resultLangSkill.Error != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create language skill"})
    return
}

c.JSON(http.StatusOK, gin.H{"langSkill": langSkill})

}

func EditLangSkill(c *gin.Context) {
    // Get profile ID and skill ID from the request parameters
    profileID, err := strconv.ParseUint(c.Param("profileID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProfileID"})
        return
    }

    // Find the profile by ID
    var profile models.Profile
    if err := initializers.DB.Preload("Skill.Lang_skill").First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Check if there are any language skills
    if len(profile.Skill) == 0 || len(profile.Skill[1].Lang_skill) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No language skills found for the profile"})
        return
    }

    // Parse the request body to get the new title for the first language skill
    var requestBody struct {
        Title string `json:"title"`
    }
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Update the title of the first language skill
    firstLangSkill := profile.Skill[1].Lang_skill[0]
    firstLangSkill.Title = requestBody.Title
    if err := initializers.DB.Save(&firstLangSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Language skill updated successfully"})
}


func DeleteLangSkill(c *gin.Context) {
    // Get profile ID and language skill ID from the request parameters
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

    // Check if the language skill exists and belongs to the correct profile's skill
    var langSkill models.Lang_skill
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&langSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Language skill not found or does not belong to the profile"})
        return
    }

    // Delete the language skill from the database
    if err := initializers.DB.Delete(&langSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Language skill deleted successfully"})
}

func GetAllLanguageSkills(c *gin.Context) {
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

    // Query language skills associated with the profile
    var languageSkills []models.Lang_skill
    if err := initializers.DB.Model(&models.Lang_skill{}).Select("title").Joins("JOIN skills ON lang_skills.skill_id = skills.id").Where("skills.profile_id = ?", profileID).Find(&languageSkills).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve language skills"})
        return
    }

    // Extract titles from language skills
    var titles []string
    for _, langSkill := range languageSkills {
        titles = append(titles, langSkill.Title)
    }

    c.JSON(http.StatusOK, gin.H{"languageSkills": titles})
}
