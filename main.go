package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("group/:gid/position", controllers.GetPosition)
	r.POST("group/:gid/position", controllers.AddPosition)
	r.DELETE("group/:gid/position/:pid", controllers.DeletePosition)
	r.PUT("group/:gid/position/:pid", controllers.EditPosition)
	r.GET("/auth", controllers.GoogleLogin)
	r.GET("/auth/callback", controllers.Googlecallback)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getusers", middleware.RequireAuth, controllers.Getusers)
	r.GET("/POSTS/:id", controllers.GetProfileWithUser)
	r.POST("/profile", controllers.ProfileCreate)
	r.PUT("/profile/:id", controllers.ProfileUpdate)
	r.GET("/profile", controllers.ProfileIndex)
	r.GET("/profile/:id", controllers.ProfileShow)
	r.DELETE("/profile/:id", controllers.ProfileDelete)
	r.POST("/POSTS", controllers.User)
	r.POST("/profiles/:profileID/exp", controllers.CreateExperience)
	r.Run()
}
