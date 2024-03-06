package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/gin-gonic/gin"
)

func MemberRoutes(r *gin.Engine) {
	// POST
	r.POST("/add-member", controllers.AddMember)
	r.POST("/remove-member", controllers.DeleteMember)
	r.POST("/group/:gid/position2", controllers.AddPositionWithSkill)
}
