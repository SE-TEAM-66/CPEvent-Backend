package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func Apply(c *gin.Context) {
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

	if err := initializers.DB.First(&models.ReqPosition{}, pid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PositionID not found!"})
		return
	}

	// Model call
	applicant, _ := c.Get("user")
	applicant, ok := applicant.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type in the context"})
		return
	}

	fmt.Println(applicant)

	var position models.ReqPosition
	if err := initializers.DB.First(&position, pid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Position not found!"})
		return
	}

	if err := initializers.DB.Model(applicant).Association("Applications").Append(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// return
	c.Status(http.StatusOK)
}
