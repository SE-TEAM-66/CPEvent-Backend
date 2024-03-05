package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.Engine) {

	// GET
	r.GET("/profile", controllers.ProfileIndex)
	r.GET("/user_profile", middleware.RequireAuth, controllers.ProfileImage)
	r.GET("/profile/:id", controllers.ProfileShow)

	// POST
	r.POST("/profile", controllers.ProfileCreate)
	r.POST("/profiles/:profileID/exp", controllers.CreateExperience)

	// PUT
	r.PUT("/profile/:id", controllers.ProfileUpdate)

	// DELETE
	r.DELETE("/profile/:id", controllers.ProfileDelete)

}
