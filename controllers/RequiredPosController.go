package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
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
	initializers.DB.Model(&group).Association("ReqPositions").Find(&all_positions)

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
