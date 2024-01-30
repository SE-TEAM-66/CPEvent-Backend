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

	r.POST("/profile", controllers.ProfileCreate)


	r.Run() 
}
