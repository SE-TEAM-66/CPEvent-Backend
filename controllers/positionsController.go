package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func GetPosition(c *gin.Context) {
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	positions, err := models.GetPosition(uint(gid))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": positions})
}

func AddPosition(c *gin.Context) {
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.AddPosition(uint(gid), position.Role)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DeletePosition(c *gin.Context) {
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = models.ValidPosition(uint(gid), uint(pid))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = models.DeletePosition(uint(pid))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
