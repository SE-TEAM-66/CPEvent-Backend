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
	var body struct{
		Gname 		string 	`json:"gname" binding:"required"`
		Owner_id 	int 	`json:"owner_id" binding:"required,gt=0"`
		Topic 		string	`json:"topic" binding:"required"`
		Description string	`json:"description" binding:"required"`
		IsHidden 	bool	`json:"is_hidden"`
		Limit_mem 	int		`json:"limit_mem" binding:"required,gt=0"`
		Cat_id 		int		`json:"cat_id" binding:"required,gt=0"`
	}
	
	// Bind and validate
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Find creator id
	var owner models.Profile
	if err := initializers.DB.First(&owner,body.Owner_id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "profile not found",
		})
		return
	}

	//Create Group
	group := models.Group{
		Gname: body.Gname,
		Owner_id: body.Owner_id,
		Topic: body.Topic,
		Description: body.Description,
		IsHidden: body.IsHidden,
		Limit_mem: body.Limit_mem,
		Cat_id: body.Cat_id, 
	}

	if err := initializers.DB.Create(&group); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create group",
		})
		return
	}

	// associate the created group to the owner id
	initializers.DB.Model(&group).Association("Profiles").Append(&owner)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": owner,
	})
}

func JoinGroup(c * gin.Context){
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
	if total_count >= int64(group.Limit_mem){
		c.JSON(http.StatusForbidden, gin.H{
			"error": "this group is full",
		})
		return			
	}

	//Check if member is duplicate in group
	var dup_count int64
	initializers.DB.Table("group_member").Where("profile_id = ?", profile.ID).Where("group_id = ?", group.ID).Count(&dup_count)
	if dup_count >= 1{
		c.JSON(http.StatusForbidden, gin.H{
			"error": "this member is already part of group",
		})
		return			
	}	

	// associate the group to the owner id
	initializers.DB.Model(&group).Association("Profiles").Append(&profile)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"count": total_count,		
		"profile": profile,
		"group": group,
	})	

}

func LeftGroup(c *gin.Context){
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
	if exist_count <= 0{
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
	if total_count <= 0{
		initializers.DB.Unscoped().Delete(&group)
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

func GroupInfoUpdate(c *gin.Context){
	// Get id from param
    gidStr := c.Param("gid")
    gid, err := strconv.Atoi(gidStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
        return
    }

	// Get data from req
	var body struct{
		Gname 		string 	`json:"gname"`
		Owner_id 	int 	`json:"owner_id" binding:"gt=0"`
		Topic 		string	`json:"topic"`
		Description string	`json:"description"`
		IsHidden 	bool	`json:"is_hidden"`
		Limit_mem 	int		`json:"limit_mem" binding:"gt=0"`
		Cat_id 		int		`json:"cat_id" binding:"gt=0"`
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
	initializers.DB.First(&group, gid)

	// Update Group
	result := initializers.DB.Model(&group).Updates(models.Group{
		Gname: body.Gname,
		Owner_id: body.Owner_id,
		Topic: body.Topic,
		Description: body.Description,
		IsHidden: body.IsHidden,
		Limit_mem: body.Limit_mem,
		Cat_id: body.Cat_id, 		
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

func GetAllGroupMembers(c *gin.Context){
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

func GetSingleGroup(c *gin.Context){
    // Get ID from param, handle potential errors using `ToInt()` for conversion
    gidStr := c.Param("gid")
    gid, err := strconv.Atoi(gidStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
        return
    }

    // Find Group, pre-populate fields to optimize query
    var group models.Group
    result := initializers.DB.Preload("ReqPositions").Preload("Profiles").First(&group, gid)

    // Handle errors gracefully
    if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return	
    }

    // Success response, optionally format data or remove sensitive information
    c.JSON(http.StatusOK, gin.H{
		"message": group,
	})
}

func GetAllGroups(c * gin.Context){
	//Get all groups
	var groups []models.Group
	result := initializers.DB.Preload("ReqPositions").Preload("Profiles").Find(&groups)

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

func GroupDelete(c *gin.Context){
	// Get data from id
	gid := c.Param("gid")

	// Find group
	var group models.Group
	if err := initializers.DB.First(&group, gid); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "target group not found",
		})
		return		
	}

	// Delete group and associated members and req_positions
	if err := initializers.DB.Select("Profiles", "ReqPositions").Unscoped().Delete(&group); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to delete group",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

