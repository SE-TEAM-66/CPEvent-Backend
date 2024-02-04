package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
	initializers.SyncDatabase()
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
	r.GET("/auth", controllers.GoogleLogin)
	r.GET("/auth/callback", controllers.Googlecallback)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getusers", middleware.RequireAuth, controllers.Getusers)
	


	r.POST("/profile", controllers.ProfileCreate)
	r.PUT("/profile/:id", controllers.ProfileUpdate)
	r.GET("/profile", controllers.ProfileIndex)
	r.GET("/profile/:id", controllers.ProfileShow)
	r.DELETE("/profile/:id", controllers.ProfileDelete)


	r.POST("/profiles/:profileID/exp", controllers.CreateExperience)

	r.POST("/profiles/:profileID/soft-skills", controllers.CreateSoftSkill)
	r.POST("/profiles/:profileID/lang-skills", controllers.CreateLangSkill)

	r.Run() 
}

