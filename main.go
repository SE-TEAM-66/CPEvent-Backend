package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("group/:gid/position", controllers.GetPosition)
	r.POST("group/:gid/position", controllers.AddPosition)
	r.DELETE("group/:gid/position/:pid", controllers.DeletePosition)
	r.PUT("group/:gid/position/:pid", controllers.EditPosition)
	r.Run()
}
