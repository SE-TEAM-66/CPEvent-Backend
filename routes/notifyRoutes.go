package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func NotifyRoutes(r *gin.Engine) {
	r.POST("/notify/new", controllers.NotifyCreate)
	r.GET("/notify/get", middleware.RequireAuth,controllers.NotifyGet)
}