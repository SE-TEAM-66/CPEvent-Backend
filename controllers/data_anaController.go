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
    if err := initializers.DB.Preload("Skill.DataAna").First(&profile, profileID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
        return
    }

    // Check if there are any DataAna skills
    if len(profile.Skill) == 0 || len(profile.Skill[2].DataAna) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No DataAna skills found for the profile"})
        return
    }

    // Parse the request body to get the new dataAna for the first DataAna skill
    var requestBody struct {
        DataAna string `json:"dataAna"`
    }
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Update the dataAna of the first DataAna skill
    firstDataAnaSkill := profile.Skill[2].DataAna[0]
    firstDataAnaSkill.DataAna = requestBody.DataAna
    if err := initializers.DB.Save(&firstDataAnaSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DataAna skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "DataAna skill updated successfully"})
}


func DeleteDataAnaSkill(c *gin.Context) {
    // Get profile ID and data analysis skill ID from the request parameters
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

    // Check if the data analysis skill exists and belongs to the correct profile's skill
    var dataAnaSkill models.DataAna
    if err := initializers.DB.Where("id = ? AND skill_id IN (SELECT id FROM skills WHERE profile_id = ?)", skillID, profileID).First(&dataAnaSkill).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data analysis skill not found or does not belong to the profile"})
        return
    }

    // Delete the data analysis skill from the database
    if err := initializers.DB.Delete(&dataAnaSkill).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data analysis skill"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Data analysis skill deleted successfully"})
}

func GetAllDataAnalysisSkills(c *gin.Context) {
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

    // Query data analysis skills associated with the profile
    var dataAnalysisSkills []models.DataAna
    if err := initializers.DB.Model(&models.DataAna{}).Select("data_ana").Joins("JOIN skills ON data_anas.skill_id = skills.id").Where("skills.profile_id = ?", profileID).Find(&dataAnalysisSkills).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data analysis skills"})
        return
    }

    // Extract data analysis from skills
    var dataAnalysis []string
    for _, skill := range dataAnalysisSkills {
        dataAnalysis = append(dataAnalysis, skill.DataAna)
    }

    c.JSON(http.StatusOK, gin.H{"dataAnalysisSkills": dataAnalysis})
}
