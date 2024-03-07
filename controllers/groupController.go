package controllers

import (
	"net/http"
	"strconv"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func GroupCreate(c *gin.Context) {
	//Get data from req body
	var body struct {
		Gname       string `json:"gname"`
		Topic       string `json:"topic" binding:"required"`
		Description string `json:"description"`
		IsHidden    bool   `json:"is_hidden"`
		Limit_mem   int    `json:"limit_mem"`
		Cat_id      int    `json:"cat_id" binding:"required,gt=0"`
	}

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

	// Fetch the profile using the user's ID
	var ownerProfile models.Profile
	result := initializers.DB.Where("user_id = ?", userModel.ID).Find(&ownerProfile)
	if result.Error != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Bind and validate
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Create Group
	group := models.Group{
		Gname:       body.Gname,
		Owner_id:    ownerProfile.ID,
		Topic:       body.Topic,
		Description: body.Description,
		IsHidden:    body.IsHidden,
		Limit_mem:   0,
		Cat_id:      body.Cat_id,
	}

	if err := initializers.DB.Create(&group); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create group",
		})
		return
	}

	owner := models.Member{
		ProfileID: ownerProfile.ID,
		GroupID:   group.ID,
		Role:      "Owner",
	}

	// associate the created group to the owner id
	initializers.DB.Model(&group).Association("Members").Append(&owner)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": group,
	})
}

func JoinGroup(c *gin.Context) {
	// Get ids from params
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Find group
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return
	}

	// Find user profile
	var profile models.Profile
	if err := initializers.DB.First(&profile, pid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target profile not found",
		})
		return
	}

	// Check if group is full
	var total_count int64
	initializers.DB.Table("group_member").Where("group_id = ?", group.ID).Count(&total_count)
	if total_count >= int64(group.Limit_mem) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "this group is full",
		})
		return
	}

	//Check if member is duplicate in group
	var dup_count int64
	initializers.DB.Table("group_member").Where("profile_id = ?", profile.ID).Where("group_id = ?", group.ID).Count(&dup_count)
	if dup_count >= 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "this member is already part of group",
		})
		return
	}

	// associate the group to the owner id
	initializers.DB.Model(&group).Association("Profiles").Append(&profile)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"count":   total_count,
		"profile": profile,
		"group":   group,
	})

}

func LeftGroup(c *gin.Context) {
	// Get ids from params
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Find group
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return
	}

	// Find user
	var profile models.Profile
	if err := initializers.DB.First(&profile, pid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target profile not found",
		})
		return
	}

	// Check if this user profile is in the group
	var exist_count int64
	initializers.DB.Table("group_member").Where("profile_id = ?", profile.ID).Where("group_id = ?", group.ID).Count(&exist_count)
	if exist_count <= 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "this profile is not the group member",
		})
		return
	}

	// Remove user profile from group association
	initializers.DB.Model(&group).Association("Profiles").Delete(&profile)

	// Check if group is empty, if yes then delete it
	var total_count int64
	initializers.DB.Table("group_member").Where("group_id = ?", group.ID).Count(&total_count)
	if total_count <= 0 {
		initializers.DB.Select("Profiles", "ReqPositions").Unscoped().Delete(&group)
		c.JSON(http.StatusOK, gin.H{
			"message": "member removed, and group deleted",
		})
		return
	}

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": "member removed",
	})
}

func GroupInfoUpdate(c *gin.Context) {
	// Get id from param
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Get data from req
	var body struct {
		Gname       string `json:"gname"`
		Topic       string `json:"topic"`
		Description string `json:"description"`
		IsHidden    bool   `json:"is_hidden"`
		Limit_mem   int    `json:"limit_mem" binding:"gt=0"`
		Cat_id      int    `json:"cat_id" binding:"gt=0"`
	}

	// Bind and validate
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Find Group from id
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return
	}

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

	// Update Group
	result := initializers.DB.Model(&group).Updates(models.Group{
		Gname:       body.Gname,
		Owner_id:    userModel.ID,
		Topic:       body.Topic,
		Description: body.Description,
		IsHidden:    body.IsHidden,
		Limit_mem:   body.Limit_mem,
		Cat_id:      body.Cat_id,
	})

	//Return on error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to update group",
		})
		return
	}

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func GetAllGroupMembers(c *gin.Context) {
	// Get gid from param
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Find group
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return
	}

	// Retrieve all user profiles from group
	var profiles []models.Profile
	initializers.DB.Model(&group).Association("Profiles").Find(&profiles)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": profiles,
	})
}

func GetSingleGroup(c *gin.Context) {
	// Get ID from param, handle potential errors using `ToInt()` for conversion
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Find Group, pre-populate fields to optimize query
	var group models.Group
	result := initializers.DB.Preload("ReqPositions").Preload("ReqPositions.Applicants").Preload("Members").Preload("Members.Profile").Preload("Members.Skills").Preload("ReqPositions.Skills").First(&group, gid)

	// Handle errors gracefully
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
			"text":  result.Error,
		})
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

	isYourGroup := userModel.Profile.ID == group.Owner_id

	// Success response, optionally format data or remove sensitive information
	c.JSON(http.StatusOK, gin.H{
		"message": group,
		"isYour":  isYourGroup,
		"profile": userModel,
	})
}

func GetAllGroups(c *gin.Context) {
	//Get all groups
	var groups []models.Group
	result := initializers.DB.Preload("ReqPositions").Preload("Members").Preload("Members.Profile").Find(&groups)

	//Return on error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to all groups",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": groups,
	})
}

func GroupDelete(c *gin.Context) {
	// Get data from id
	gidStr := c.Param("gid")
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	// Find group
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return
	}

	// Delete group and associated members and req_positions
	if err := initializers.DB.Select("Members", "ReqPositions").Unscoped().Delete(&group); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to delete group",
		})
		return
	}

	// Response w
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func GetUserGroups(c *gin.Context) {
	// Get user ID from the context
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

	// Retrieve all groups that the user is a member of
	var userGroups []models.Group
	initializers.DB.Find(&userGroups, "id IN (SELECT group_id FROM members WHERE profile_id = ?)", userProfile.ID)

	// Return on success
	c.JSON(http.StatusOK, gin.H{
		"message": userGroups,
	})
}
