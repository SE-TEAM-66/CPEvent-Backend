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
	if err := initializers.DB.Where("id = ?", gid).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	var positions []models.Position
	if err := initializers.DB.Where("group_id = ?", gid).Find(&positions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// return
	c.JSON(http.StatusOK, gin.H{"data": positions})
}

func AddPosition(c *gin.Context) {
	// Get request
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GroupID!"})
		return
	}
	var reqPos models.Position
	if err := c.ShouldBindJSON(&reqPos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model Call
	var group models.Group
	if err := initializers.DB.Where("id = ?", gid).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID not found!"})
		return
	}

	newPos := models.Position{Role: reqPos.Role, GroupID: uint(gid)}
	initializers.DB.Model(&group).Association("Positions").Append(&newPos)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PrositionID!"})
		return
	}

	if err := validPosition(uint(gid), uint(pid)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GroupID or PrositionID incorrect!"})
		return
	}

	// Model Call
	var position models.Position

	if err := initializers.DB.Delete(&position, pid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return
	c.Status(http.StatusOK)
}

func validPosition(gid uint, pid uint) error {
	var positions []models.Position
	result := initializers.DB.Where("group_id = ?", gid).First(&positions, pid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
