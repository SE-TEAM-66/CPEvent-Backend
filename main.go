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
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getusers", middleware.RequireAuth, controllers.Getusers)
	r.GET("/group/:gid", controllers.GetSingleGroup)
	r.GET("/all-groups", controllers.GetAllGroups)
	r.POST("/new-group", controllers.GroupCreate)
	r.POST("/group/:gid/add/:uid", controllers.JoinGroup)
	r.PUT("/set-group/:gid", controllers.GroupInfoUpdate)
	r.DELETE("/rm-group/:gid", controllers.GroupDelete)
	r.DELETE("/group/:gid/rm/:uid", controllers.LeftGroup)
	r.Run()
}
