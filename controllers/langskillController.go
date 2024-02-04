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
