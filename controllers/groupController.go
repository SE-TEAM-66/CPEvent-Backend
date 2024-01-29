package controllers

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
)

func GroupCreate(c *gin.Context) {
	//Get data from req body
	var body struct{
		Gname string 
		Owner_id int
		Topic string
		Description string
		IsHidden bool
		Limit_mem int
		Cat_id int
	}
	c.Bind(&body)

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
	result := initializers.DB.Create(&group)

	//Return on error
	if result.Error != nil {
		c.Status(400)
		return
	}

	//Return on Success
	c.JSON(200, gin.H{
		"group": group,
	})
}

func GroupUpdate(c *gin.Context){
	// Get id from param
	id := c.Param("gid")

	// Get data from req
	var body struct {
		Gname string 
		Owner_id int
		Topic string
		Description string
		IsHidden bool
		Limit_mem int
		Cat_id int	
	}
	c.Bind(&body)

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
		c.Status(400)
		return
	}

	//Return on Success
	c.JSON(200, gin.H{
		"group": group,
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
		c.Status(400)
		return
	}

	//Response 
	c.JSON(200, gin.H{
		"group": group,
	})

}

func GetAllGroups(c * gin.Context){
	//Get all groups
	var groups []models.Group
	result := initializers.DB.Find(&groups)

	//Return on error
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Response
	c.JSON(200, gin.H{
		"groups": groups,
	})
}

func GroupDelete(c *gin.Context){
	// Get data from id
	id := c.Param("gid")

	// Delete group
	result := initializers.DB.Delete(&models.Group{},id)

	//Return on error
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Response
	c.Status(200)
}
