package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

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

