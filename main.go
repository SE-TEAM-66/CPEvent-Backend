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
	r.GET("/auth", controllers.GoogleLogin)
	r.GET("/auth/callback", controllers.Googlecallback)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.POST("/logoutg", controllers.Logoutg)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getusers", middleware.RequireAuth, controllers.Getusers)
	r.Run()
}
