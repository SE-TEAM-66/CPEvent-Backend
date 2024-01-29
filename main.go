package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func init(){
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	
	r.POST("/new-group", controllers.GroupCreate)
	r.PUT("/group/:gid", controllers.GroupUpdate)

	r.Run() 
}