package main

import (
	"fmt"
	"os"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/SE-TEAM-66/CPEvent-Backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
	initializers.SyncDatabase()

	initializers.DB.AutoMigrate(&models.Group{})
	initializers.DB.AutoMigrate(&models.ReqPosition{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Profile{})
	initializers.DB.AutoMigrate(&models.Exp{})
	initializers.DB.AutoMigrate(&models.Skill{})
	initializers.DB.AutoMigrate(&models.DataAna{})
	initializers.DB.AutoMigrate(&models.Lang_skill{})
	initializers.DB.AutoMigrate(&models.Soft_skill{})
	initializers.DB.AutoMigrate(&models.GroupSkill{})
	initializers.DB.AutoMigrate(&models.Member{})
	initializers.DB.AutoMigrate(&models.Notify{})
	initializers.DB.AutoMigrate(&models.Event{})
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{os.Getenv("ALLOW_ORIGINS_URL")}
	fmt.Println(os.Getenv("ALLOW_ORIGINS_URL"))
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	routes.GroupRoutes(r)
	routes.AuthRoutes(r)
	routes.ProfileRoutes(r)
	routes.EventRoutes(r)
	routes.MemberRoutes(r)
	routes.NotifyRoutes(r)

	r.Run()
}
