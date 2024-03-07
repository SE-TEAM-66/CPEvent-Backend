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

func AcceptApplicant(c *gin.Context) {
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

	applicantID, err := strconv.ParseUint(c.Param("applicantID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ApplicantID!"})
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

	// Check if the user is a member of the group
	var existingMember models.Member
	result = initializers.DB.Where("group_id = ? AND profile_id = ?", gid, userProfile.ID).First(&existingMember)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	// Check if the user is the owner of the group
	var group models.Group
	result = initializers.DB.Where("id = ? AND owner_id = ?", gid, userProfile.ID).First(&group)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have the permission to accept applicants"})
		return
	}

	// Check if the position exists and is associated with the group
	var reqPosition models.ReqPosition
	result = initializers.DB.Where("group_id = ? AND id = ?", gid, pid).First(&reqPosition)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID or GroupID"})
		return
	}

	// Check if the applicant exists and is in the Applicants list for the position
	var applicantProfile models.Profile
	result = initializers.DB.Where("id = ?", applicantID).First(&applicantProfile)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ApplicantID"})
		return
	}

	// Remove the applicant from the Applicants list
	if err := initializers.DB.Model(&reqPosition).Association("Applicants").Delete(&applicantProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add the applicant as a member of the group
	newMember := models.Member{
		Role:      "member", // You can set the role as needed
		ProfileID: applicantProfile.ID,
		GroupID:   uint(gid),
	}

	if err := initializers.DB.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Applicant accepted into the group"})
}

func RejectApplicant(c *gin.Context) {
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

	applicantID, err := strconv.ParseUint(c.Param("applicantID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ApplicantID!"})
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

	// Check if the user is a member of the group
	var existingMember models.Member
	result = initializers.DB.Where("group_id = ? AND profile_id = ?", gid, userProfile.ID).First(&existingMember)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	// Check if the user is the owner of the group
	var group models.Group
	result = initializers.DB.Where("id = ? AND owner_id = ?", gid, userProfile.ID).First(&group)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have the permission to reject applicants"})
		return
	}

	// Check if the position exists and is associated with the group
	var reqPosition models.ReqPosition
	result = initializers.DB.Where("group_id = ? AND id = ?", gid, pid).First(&reqPosition)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID or GroupID"})
		return
	}

	// Check if the applicant exists and is in the Applicants list for the position
	var applicantProfile models.Profile
	result = initializers.DB.Where("id = ?", applicantID).First(&applicantProfile)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ApplicantID"})
		return
	}

	// Remove the applicant from the Applicants list
	if err := initializers.DB.Model(&reqPosition).Association("Applicants").Delete(&applicantProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Applicant rejected from the group"})
}

func CancelApplication(c *gin.Context) {
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

	// Check if the position exists and is associated with the group
	var reqPosition models.ReqPosition
	result = initializers.DB.Where("group_id = ? AND id = ?", gid, pid).First(&reqPosition)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PositionID or GroupID"})
		return
	}

	// Check if the user's profile is in the Applicants list for the position
	var applicantProfile models.Profile
	if err := initializers.DB.Model(&reqPosition).Association("Applicants").Find(&applicantProfile, "id = ?", userProfile.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove the applicant from the Applicants list
	if err := initializers.DB.Model(&reqPosition).Association("Applicants").Delete(&applicantProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application canceled successfully"})
}
