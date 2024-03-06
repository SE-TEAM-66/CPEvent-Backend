package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}

func main() {
	groupDB()
}

func groupDB() {
	initializers.DB.AutoMigrate(&models.Group{})
	initializers.DB.AutoMigrate(&models.ReqPosition{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Profile{})
	initializers.DB.AutoMigrate(&models.Exp{})
	initializers.DB.AutoMigrate(&models.GroupSkill{})
	initializers.DB.AutoMigrate(&models.Member{})
	initializers.DB.AutoMigrate(&models.Notify{})
}
