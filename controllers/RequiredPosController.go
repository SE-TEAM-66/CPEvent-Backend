package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosition(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}

	// Model Call
	var group models.Group
	if err := initializers.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	var all_positions []models.ReqPosition
	initializers.DB.Preload("Skills").Find(&all_positions, "req_positions.group_id = ?", gid)

	user, _ := c.Get("user")
	var user_positions []models.ReqPosition
	initializers.DB.Model(user).Association("Applications").Find(&user_positions)

	type Position struct {
		Position models.ReqPosition
		Applied  bool
	}

	positions := []Position{}

	for _, pos := range all_positions {
		pass := false
		for _, user_pos := range user_positions {
			if pos.ID == user_pos.ID {
				position := Position{Position: pos, Applied: true}
				positions = append(positions, position)
				pass = true
				break
			}
		}
		if !pass {
			position := Position{Position: pos, Applied: false}
			positions = append(positions, position)
		}
	}

	// return
	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

func AddPosition(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}
	var newPos models.ReqPosition
	if err := c.ShouldBindJSON(&newPos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model Call
	var group models.Group
	if err := initializers.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	initializers.DB.Model(&group).Association("ReqPositions").Append(&newPos)

	// return
	c.JSON(http.StatusOK, gin.H{"data": newPos})
}

func DeletePosition(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID!"})
		return
	}

	if err := initializers.DB.First(&models.Group{}, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	// Model Call
	var targetPos models.ReqPosition
	if err := initializers.DB.First(&targetPos, pid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PositionID not found!"})
		return
	}

	if err := validPosition(uint(gid), uint(pid)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID and PositionID inconsistent!"})
		return
	}

	if err := initializers.DB.Unscoped().Delete(&targetPos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return
	c.Status(http.StatusOK)
}

func EditPosition(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID!"})
		return
	}

	if err := initializers.DB.First(&models.Group{}, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	var reqPos models.ReqPosition
	if err := c.ShouldBindJSON(&reqPos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model Call
	var editedPos models.ReqPosition
	if err := initializers.DB.First(&editedPos, pid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PositionID not found!"})
		return
	}

	if err := validPosition(uint(gid), uint(pid)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID and PositionID inconsistent!"})
		return
	}

	if err := initializers.DB.Model(&editedPos).Updates(
		models.ReqPosition{
			Role: reqPos.Role,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return
	c.Status(http.StatusOK)
}

func validPosition(gid uint, pid uint) error {
	var positions []models.ReqPosition
	result := initializers.DB.Where("group_id = ?", gid).First(&positions, pid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddPositionWithSkill(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}

	var newPos models.ReqPosition
	if err := c.ShouldBindJSON(&newPos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if skills are provided
	if len(newPos.Skills) > 0 {
		// Validate and associate skills with the new position
		for i, skill := range newPos.Skills {
			if err := validateSkill(&skill); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "skill": newPos})
				return
			}

			// Check if the skill already exists
			var existingSkill models.GroupSkill
			result := initializers.DB.Where("name = ?", skill.Name).First(&existingSkill)
			if result.Error != nil {
				if result.Error != gorm.ErrRecordNotFound {
					// Other database error
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}
			}

			// If the skill exists, use the existing one; otherwise, create a new skill
			if result.Error == nil {
				newPos.Skills[i] = existingSkill
			} else {
				// Skill doesn't exist, create a new one
				if err := initializers.DB.Create(&skill).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				newPos.Skills[i] = skill
			}
		}
	}

	// Model Call
	var group models.Group
	if err := initializers.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	// Create new position
	newPos.GroupID = group.ID
	if err := initializers.DB.Create(&newPos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return
	c.JSON(http.StatusOK, gin.H{"data": newPos})
}

func validateSkill(skill *models.GroupSkill) error {
	// Implement your skill validation logic here
	// For example, you can check if the skill name is not empty
	if skill.Name == "" {
		return errors.New("skill name cannot be empty")
	}
	return nil
}

func ApplyPosition(c *gin.Context) {
	// Get request parameters
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}

	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID!"})
		return
	}

	// Check if the user making the request is the owner of the group
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in the context"})
		return
	}

	// Type assertion
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type in the context"})
		return
	}

	// Find user's profile
	var userProfile models.Profile
	result := initializers.DB.Where("user_id = ?", userModel.ID).First(&userProfile)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Check if the group exists
	var existingGroup models.Group
	if err := initializers.DB.First(&existingGroup, gid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	// Check if the position exists in the group
	var existingPosition models.ReqPosition
	if err := initializers.DB.Where("id = ? AND group_id = ?", pid, gid).First(&existingPosition).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PositionID not found in the specified group!"})
		return
	}

	// Check if the user is already a member of the group
	var existingMember models.Member
	result = initializers.DB.Where("group_id = ? AND profile_id = ?", gid, userProfile.ID).First(&existingMember)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already a member of this group!"})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		// Other database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Check if the user has already applied for the position
	var existingApplication models.ReqPosition
	result = initializers.DB.Preload("Applicants").Where("group_id = ? AND ID = ?", gid, pid).First(&existingApplication)
	if result.Error == nil {
		// Check if the user's profile is already in the Applicants list
		for _, applicant := range existingApplication.Applicants {
			if applicant.ID == userProfile.ID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "You have already applied for this position!"})
				return
			}
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		// Other database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Add the user's profile to the Applicants list
	if err := initializers.DB.Model(&existingApplication).Association("Applicants").Append(&userProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully!"})
}
