package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	// GET
	r.GET("/auth", controllers.GoogleLogin)
	r.GET("/auth/callback", controllers.Googlecallback)
	r.GET("/vali_date", middleware.RequireAuth, controllers.Validate)
	r.GET("/get_users", middleware.RequireAuth, controllers.Getusers)
	r.GET("/log_out", middleware.RequireAuth, controllers.Logout)
	// POST
	r.POST("/sign_up", controllers.Signup)
	r.POST("/log_in", controllers.Login)

}
