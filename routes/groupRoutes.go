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
	r.GET("/group/all", controllers.GetAllGroups)

	// POST
	r.POST("/group/new", controllers.GroupCreate)
	r.POST("/group/:gid/add/:pid", controllers.JoinGroup)

	// PUT
	r.PUT("/group/set/:gid", controllers.GroupInfoUpdate)

	// DELETE
	r.DELETE("/group/:gid/rm/:pid", controllers.LeftGroup)
	r.DELETE("/group/del/:gid", controllers.GroupDelete)

	// Position
	// GET
	r.GET("/group/:gid/position", middleware.RequireAuth, controllers.GetPosition)

	// POST
	r.POST("/group/:gid/position", controllers.AddPosition)

	// PUT
	r.PUT("/group/:gid/position/:pid", controllers.EditPosition)

	// DELETE
	r.DELETE("/group/:gid/position/:pid", controllers.DeletePosition)

	// apply position
	r.POST("group/:gid/position/:pid", middleware.RequireAuth, controllers.Apply)

}
