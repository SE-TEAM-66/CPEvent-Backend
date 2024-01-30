package controllers

import (

	"net/http"

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
	var owner models.User
	if err := initializers.DB.First(&owner,body.Owner_id); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
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
	initializers.DB.Model(&group).Association("Users").Append(&owner)

	//Return on Success
	c.JSON(http.StatusOK, gin.H{
		"message": owner,
	})
}

func GroupUpdate(c *gin.Context){
	// Get id from param
	id := c.Param("gid")

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
	initializers.DB.First(&group, id)

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

func GetSingleGroup(c *gin.Context){
	//Get id from param
	id := c.Param("gid")

	// Find Group
	var group models.Group
	result := initializers.DB.First(&group, id)

	//Return on error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get group from id",
		})
		return
	}

	//Response 
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

}

func GetAllGroups(c * gin.Context){
	//Get all groups
	var groups []models.Group
	result := initializers.DB.Find(&groups)

	//Return on error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to all groups",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func GroupDelete(c *gin.Context){
	// Get data from id
	id := c.Param("gid")

	// Delete group
	result := initializers.DB.Delete(&models.Group{},id)

	//Return on error
	if result.Error != nil {
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
