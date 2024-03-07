package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func MemberRoutes(r *gin.Engine) {
	// POST
	r.POST("/add-member", controllers.AddMember)
	r.POST("/remove-member", controllers.DeleteMember)
	r.POST("/group/:gid/position2", controllers.AddPositionWithSkill)
	r.POST("/accept-applicant/:gid/:pid/:applicantID", middleware.RequireAuth, controllers.AcceptApplicant)
	r.POST("/reject-applicant/:gid/:pid/:applicantID", middleware.RequireAuth, controllers.RejectApplicant)
	r.POST("/cancel-applicant/:gid/:pid", middleware.RequireAuth, controllers.CancelApplication)
}
