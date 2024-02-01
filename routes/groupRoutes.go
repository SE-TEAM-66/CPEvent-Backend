package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func GroupRoutes(r *gin.Engine) {
	// Group
		// GET
	r.GET("/group/:gid/all-members", controllers.GetAllGroupMembers)
	r.GET("/group/:gid", controllers.GetSingleGroup)	
	r.GET("/all-groups", controllers.GetAllGroups)

		// POST
	r.POST("/new-group", controllers.GroupCreate)
	r.POST("/group/:gid/add/:uid", controllers.JoinGroup)

		// PUT
	r.PUT("/set-group/:gid", controllers.GroupInfoUpdate)

		// DELETE
	r.DELETE("/group/:gid/rm/:uid", controllers.LeftGroup)	
	r.DELETE("/rm-group/:gid", controllers.GroupDelete)

	// Position
		// GET
	r.GET("/group/:gid/position", middleware.RequireAuth, controllers.GetPosition)

		// POST
	r.POST("/group/:gid/position", controllers.AddPosition)

		// PUT
	r.PUT("/group/:gid/position/:pid", controllers.EditPosition)

		// DELETE
	r.DELETE("/group/:gid/position/:pid", controllers.DeletePosition)	

}