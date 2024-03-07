package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.Engine) {
	// GET
	r.GET("/event/all",controllers.GetAllEvents)
	r.GET("/event/:eid", controllers.GetSingleEvent)

	// POST
	r.POST("/event/new",controllers.EventCreate)
}